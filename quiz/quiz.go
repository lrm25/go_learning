package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	cmdLineArgs := os.Args[1:]
	if 1 < len(cmdLineArgs) {
		fmt.Println("Too many arguments.")
		fmt.Println("Format: \"quiz.exe <optional csv file name>\"")
		os.Exit(1)
	}

	var filename string
	if 1 == len(cmdLineArgs) {
		filename = cmdLineArgs[0]
	} else {
		filename = "problems.csv"
	}

	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		defer csvFile.Close()

		csvReader := csv.NewReader(bufio.NewReader(csvFile))
		right := 0
		wrong := 0

		var inReader *bufio.Reader = bufio.NewReader(os.Stdin)
		fmt.Printf("%T\n", inReader)

		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println(err.Error())
				continue
			}

			q := Question{}
			err = q.LoadQuestion(line)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			q.Ask(&right, &wrong, inReader)
		}
		fmt.Printf("%d correct, %d incorrect\n", right, wrong)
	}
}

/* question */

type Question struct {
	question string
	answer   string
}

func (q Question) Ask(right *int, wrong *int, reader *bufio.Reader) {
	fmt.Printf("%s ", q.question)
	userAnswer, _ := reader.ReadString('\n')
	userAnswer = strings.TrimSuffix(userAnswer, "\n")
	userAnswer = strings.TrimSuffix(userAnswer, "\r")
	if userAnswer == q.answer {
		*right += 1
		fmt.Println("Correct")
	} else {
		*wrong += 1
		fmt.Println("Incorrect")
	}
}

func (q *Question) LoadQuestion(qAndA []string) error {

	var err error = nil

	if len(qAndA) != 2 {
		err = errors.New(fmt.Sprintf("Invalid question format for \"%s\"", qAndA))
		return err
	}
	q.question = qAndA[0]
	q.answer = qAndA[1]
	return err
}
