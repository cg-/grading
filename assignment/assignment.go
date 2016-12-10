/*
 *  An object representing an Assignment to be graded.
 */

package assignment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"log"

	"path"

	"os/exec"

	"github.com/cg-/grading/common"
)

var Debug *common.DebugLogger

// Assignment represents an assignment to be graded.
type Assignment struct {
	Name      string
	Questions []Question
}

// NewAssignment is a constructor for the Assignment.
func NewAssignment() *Assignment {
	return &Assignment{}
}

// NewAssignmentFromSpec generates an Assignment from a spec file (json)
func NewAssignmentFromSpec(spec string) (*Assignment, error) {
	toReturn := Assignment{}
	err := json.Unmarshal([]byte(spec), &toReturn)
	if err != nil {
		return nil, err
	}
	return &toReturn, nil
}

func GenerateDefaultSpec(path string) error {
	buf, err := json.Marshal(GetDefaultAssignment())
	if err != nil {
		return err
	}
	var buf2 bytes.Buffer
	err = json.Indent(&buf2, buf, "", "   ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, buf2.Bytes(), 0755)
	if err != nil {
		return err
	}
	return nil
}

// Grade will take a folder from the downloaded eCommons file as an argument and
// will grade the assignment.
func (a *Assignment) Grade(inputFolder os.File, outputFolder os.File, cpPath string) {
	Debug.Print("Started grading assignment: [" + a.String() + "]")
	Debug.Print("Looking for checkpoint...")
	cp := NewCheckpoint()
	if _, err := os.Stat(cpPath); err == nil {
		Debug.Print("Checkpoint found. Loading it up.")
		b, err := ioutil.ReadFile(cpPath)
		if err != nil {
			log.Fatalf("Trouble reading checkpoint: " + err.Error())
		}
		cp, err = NewCheckpointFromString(string(b))
		if err != nil {
			log.Fatalf("Trouble creating checkpoint from string: " + err.Error())
		}
	} else {
		Debug.Print("No checkpoint found.")
	}

	assignmentFileInfo, err := inputFolder.Readdir(0)
	if err != nil {
		log.Fatalf("Trouble reading the assignment directory: " + err.Error())
	}
	for f := range assignmentFileInfo {
		assignmentDir, err := os.Open(inputFolder.Name() + "//" + assignmentFileInfo[f].Name())
		defer assignmentDir.Close()

		if err != nil {
			log.Fatalf("Trouble reading the assignment directory: " + err.Error())
		}
		stat, err := assignmentDir.Stat()
		if err != nil {
			log.Fatalf("Trouble reading the assignment directory: " + err.Error())
		}

		err = os.MkdirAll(outputFolder.Name()+"//"+assignmentFileInfo[f].Name(), 0755)
		if err != nil {
			log.Fatalf("Trouble creating the output directory: " + err.Error())
		}

		if !stat.IsDir() {
			continue
		}
		studentDirFileInfo, err := assignmentDir.Readdir(0)
		if err != nil {
			log.Fatalf("Trouble reading the student directory: " + err.Error())
		}
		for g := range studentDirFileInfo {
			studentDir, err := os.Open(inputFolder.Name() + "//" + assignmentFileInfo[f].Name() + "//" + studentDirFileInfo[g].Name())
			defer studentDir.Close()
			if err != nil {
				log.Fatalf("Trouble reading the student directory: " + err.Error())
			}
			if !cp.Exists(studentDirFileInfo[g].Name()) {
				if !studentDirFileInfo[g].IsDir() {
					continue
				}
				grade, comments, err := a.gradeStudent(studentDir)
				if err != nil {
					fmt.Printf("Error grading %s: %s", studentDir.Name(), err.Error())
					continue
				}
				err = os.MkdirAll(outputFolder.Name()+"//"+assignmentFileInfo[f].Name()+"//"+studentDirFileInfo[g].Name(), 0755)
				if err != nil {
					log.Fatalf("error making output folder: " + err.Error())
				}
				cmtPath := outputFolder.Name() + "//" + assignmentFileInfo[f].Name() + "//" + studentDirFileInfo[g].Name() + "//" + "comments.txt"
				err = ioutil.WriteFile(cmtPath, []byte(comments), 0755)
				if err != nil {
					log.Fatalf("error writing comment file: " + err.Error())
				}
				cp.Add(studentDirFileInfo[g].Name(), grade, comments)
				cp.Save(cpPath)
				cp.CSV(outputFolder.Name() + "//" + "scores.csv")
			}
		}
	}
}

func (a *Assignment) gradeStudent(studentFolder *os.File) (int, string, error) {
	subFolder, err := os.Open(studentFolder.Name() + "//" + "Submission attachment(s)")
	if err != nil {
		log.Fatalf("Couldn't open student submission folder: " + err.Error())
	}
	files, err := subFolder.Readdirnames(0)
	if err != nil {
		log.Fatalf("Couldn't get file names from student submission folder: " + err.Error())
	}
	possibleFiles := []string{}
	allFiles := []string{}
	toGradePath := ""
	for i := range files {
		allFiles = append(allFiles, files[i])
		if path.Ext(files[i]) == ".pdf" {
			possibleFiles = append(possibleFiles, files[i])
		}
	}
	if len(possibleFiles) == 1 {
		toGradePath = possibleFiles[0]
	} else if len(possibleFiles) > 1 {
		return 0, "", fmt.Errorf("multiple pdf files")
	} else {
		return 0, "", fmt.Errorf("no pdf file")
	}
	toGradeFile, err := os.Open(subFolder.Name() + "//" + toGradePath)
	defer toGradeFile.Close()
	if err != nil {
		return 0, "", fmt.Errorf("couldn't open file to grade [%s] ", err.Error())
	}
	if runtime.GOOS == "windows" {
		cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", toGradeFile.Name())
		err = cmd.Start()
	} else {
		cmd := exec.Command("evince", toGradeFile.Name())
		err = cmd.Start()
	}
	if err != nil {
		return 0, "", fmt.Errorf("couldn't open file to grade (cmd failed) [%s] ", err.Error())
	}

	totalScore := 0
	comments := ""
	for i := range a.Questions {
		qVal, qCmt := a.Questions[i].Ask()
		comments += fmt.Sprintf("Question %d: %s[%d/%d] (Correct Answer: %s)\r\n", i+1, qCmt, qVal, a.Questions[i].Value, a.Questions[i].Answer)
		totalScore += qVal
	}
	comments += fmt.Sprintf("\r\nTOTAL: %d\r\n", totalScore)

	return totalScore, comments, nil
}

func (a *Assignment) String() string {
	return a.Name
}
