package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Results of questions (correct or incorrect)
type Results struct {
	correct    int
	incorrect  int
	unanswered int
	mutex      sync.Mutex
}

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
		var results Results

		var inReader = bufio.NewReader(os.Stdin)
		fmt.Printf("Press any key to continue: ")
		inReader.ReadString('\n')
		go Countdown(&results)

		lines, err := csvReader.ReadAll()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		results.mutex.Lock()
		results.unanswered = len(lines)
		results.mutex.Unlock()

		for _, line := range lines {
			q := Question{}
			err = q.LoadQuestion(line)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			q.Ask(&results, inReader)
		}
		results.Print()
	}
}

func (results *Results) Print() {
	results.mutex.Lock()
	fmt.Printf("%d correct, %d incorrect\n", results.correct, results.incorrect+results.unanswered)
	results.mutex.Unlock()
}

func Countdown(results *Results) {

	timer := time.NewTimer(30 * time.Second)
	<-timer.C

	fmt.Println()
	fmt.Println("Time's up!")
	results.Print()
	os.Exit(0)
}

/* question */

// Question struct
type Question struct {
	question string
	answer   string
}

// Ask question to command line, and retrieve user's answer
func (q Question) Ask(results *Results, reader *bufio.Reader) {
	fmt.Printf("%s ", q.question)
	userAnswer, _ := reader.ReadString('\n')
	userAnswer = strings.TrimSuffix(userAnswer, "\n")
	userAnswer = strings.TrimSuffix(userAnswer, "\r")
	results.mutex.Lock()
	if userAnswer == q.answer {
		results.correct++
		fmt.Println("Correct")
	} else {
		results.incorrect++
		fmt.Println("Incorrect")
	}
	results.unanswered--
	results.mutex.Unlock()
}

// LoadQuestion (split question CSV data into question and answer)
func (q *Question) LoadQuestion(qAndA []string) error {

	var err error

	if len(qAndA) != 2 {
		err = fmt.Errorf("Invalid question format for \"%s\"", qAndA)
		return err
	}
	q.question = qAndA[0]
	q.answer = qAndA[1]
	return err
}
