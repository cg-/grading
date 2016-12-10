/*
 *  Checkpoint is a structure to hold onto our current progress grading.
 */

package assignment

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Checkpoint struct {
	Completed  []string
	ScoreMap   map[string]int
	CommentMap map[string]string
}

func NewCheckpoint() *Checkpoint {
	return &Checkpoint{
		ScoreMap:   make(map[string]int),
		CommentMap: make(map[string]string),
	}
}

func NewCheckpointFromString(s string) (*Checkpoint, error) {
	toReturn := &Checkpoint{}
	err := json.Unmarshal([]byte(s), toReturn)
	if err != nil {
		return nil, err
	}
	return toReturn, nil
}

func (c *Checkpoint) CSV(csvPath string) {
	csvFile, err := os.Create(csvPath)
	defer csvFile.Close()
	if err != nil {
		log.Fatalf("Trouble creating the csv file: " + err.Error())
	}
	csvWriter := csv.NewWriter(csvFile)
	for i := range c.Completed {
		csvWriter.Write([]string{c.Completed[i], strconv.Itoa(c.ScoreMap[c.Completed[i]])})
	}
	csvWriter.Flush()
}

func (c *Checkpoint) String() string {
	toReturn := ""
	for i := range c.Completed {
		toReturn += c.Completed[i]
		if i != len(c.Completed)-1 {
			toReturn += ", "
		}
	}
	return toReturn
}

func (c *Checkpoint) Save(savePath string) error {
	buf, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = ioutil.WriteFile(savePath, buf, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (c *Checkpoint) Add(name string, score int, comments string) {
	c.Completed = append(c.Completed, name)
	c.ScoreMap[name] = score
	c.CommentMap[name] = comments
}

func (c *Checkpoint) Exists(s string) bool {
	for item := range c.Completed {
		if c.Completed[item] == s {
			return true
		}
	}
	return false
}
