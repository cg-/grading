/*
 *  A question on an assignment
 */

package assignment

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Question struct {
	Question     string
	Answer       string
	Value        int
	CommonErrors []CommonError
}

type CommonError struct {
	Error     string
	Deduction int
}

func (c *CommonError) String() string {
	toReturn := ""
	toReturn += fmt.Sprintf("%s [-%d points]", c.Error, c.Deduction)
	return toReturn
}

func (q *Question) String() string {
	toReturn := ""
	toReturn += fmt.Sprintf("Question: %s [%d points]\n", q.Question, q.Value)
	toReturn += fmt.Sprintf("Answer: %s\n", q.Answer)
	for i := range q.CommonErrors {
		toReturn += fmt.Sprintf("Common Error %d: %s\n", i+1, q.CommonErrors[i].String())
	}
	return toReturn
}

func (q *Question) Ask() (int, string) {
	fmt.Printf("[%d Points] %s\n", q.Value, q.Question)
	fmt.Printf("Answer: %s\n", q.Answer)
	inputInt := len(q.CommonErrors) + 2
	for inputInt > len(q.CommonErrors)+1 || inputInt < 0 {
		fmt.Printf("0) Full credit.\n")
		i := 1
		for err := range q.CommonErrors {
			fmt.Printf("%d) %s\n", i, q.CommonErrors[err].String())
			i++
		}
		fmt.Printf("%d) Other\n", len(q.CommonErrors)+1)
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf(">")
		inputStr, err := reader.ReadString('\n')
		inputStr = strings.TrimSpace(inputStr)
		if err != nil {
			Debug.Print("Trouble reading the input string... " + err.Error())
			fmt.Println("Invalid input.")
			continue
		}

		inputInt, err = strconv.Atoi(inputStr)
		if err != nil {
			Debug.Print("Trouble converting the input string... " + err.Error())
			fmt.Println("Invalid input.")
			continue
		}
	}
	if inputInt == 0 {
		return q.Value, "Full credit."
	}
	if inputInt == len(q.CommonErrors)+1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Comment> ")
		cmt, err := reader.ReadString('\n')
		cmt = strings.TrimSpace(cmt)
		if err != nil {
			log.Fatalf("Couldn't read input.")
		}
		fmt.Printf("Deduction> ")
		ded, err := reader.ReadString('\n')
		ded = strings.TrimSpace(ded)
		if err != nil {
			log.Fatalf("Couldn't read input.")
		}
		dedInt, err := strconv.Atoi(ded)
		if err != nil {
			Debug.Print("Trouble converting the input string... " + err.Error())
			log.Fatalf("Couldn't read input.")
		}
		return (q.Value - dedInt), cmt
	}
	return (q.Value - q.CommonErrors[inputInt-1].Deduction), q.CommonErrors[inputInt-1].Error
}
