/*
 *  Grading Helper
 */

package main

import (
	"flag"
	"log"
	"os"

	"io/ioutil"

	"github.com/cg-/grading/assignment"
	"github.com/cg-/grading/common"
)

var inputPath string
var outputPath string
var cpPath string
var specString string
var debugging bool
var Debug = common.NewDebugLogger()

func parseArgs() {
	inputFlag := flag.String("input", "bulk_download", "the input directory")
	outputFlag := flag.String("output", "output", "the output directory")
	cpFlag := flag.String("checkpoint", "checkpoint", "the checkpoint file")
	specFlag := flag.String("spec", "default.spec", "the assignment spec file")
	debuggingFlag := flag.Bool("debug", false, "show debugging messages")
	genDefaultFlag := flag.Bool("default-spec", false, "generate a default spec file")
	flag.Parse()

	if *genDefaultFlag {
		assignment.GenerateDefaultSpec("default.spec")
		os.Exit(0)
	}

	buf, err := ioutil.ReadFile(*specFlag)
	if err != nil {
		log.Fatal(err)
	}
	specString = string(buf)

	inputPath = *inputFlag
	outputPath = *outputFlag
	cpPath = *cpFlag

	if *debuggingFlag {
		Debug.Enable()
	} else {
		Debug.Disable()
	}
}

func main() {
	// Inject Debug Logger to other packages that need it.
	assignment.Debug = Debug

	// Parse arguments
	parseArgs()

	assignment, err := assignment.NewAssignmentFromSpec(specString)
	if err != nil {
		log.Fatal("Error generating assignment from spec file: " + err.Error())
	}

	inFile, err := os.Open(inputPath)
	defer inFile.Close()
	if err != nil {
		log.Fatal("Error opening file (" + inputPath + "): " + err.Error())
	}
	stat, err := inFile.Stat()
	if err != nil {
		log.Fatal("Error running stat on file(" + inFile.Name() + "): " + err.Error())
	}
	if !stat.Mode().IsDir() {
		log.Fatal("Input wasn't a directory.")
	}

	outFile, err := os.Open(outputPath)
	defer outFile.Close()
	if os.IsNotExist(err) {
		err2 := os.MkdirAll(outputPath, 0755)
		if err2 != nil {
			log.Fatal("Error creating output folder: " + err.Error())
		}
		outFile, err = os.Open(outputPath)
		if err != nil {
			log.Fatal("Error opening file(" + outputPath + "): " + err.Error())
		}
	} else if err != nil {
		log.Fatal("Error opening file(" + outputPath + "): " + err.Error())
	} else {
		Debug.Print("Warning: Output folder already exists.")
	}
	stat, err = outFile.Stat()

	if err != nil {
		log.Fatal("Error running stat on file: " + err.Error())
	} else if !stat.Mode().IsDir() {
		log.Fatal("Input wasn't a directory.")
	}

	assignment.Grade(*inFile, *outFile, cpPath)
}
