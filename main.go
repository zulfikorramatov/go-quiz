package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type Problem struct {
	question string
	answer   string
}

func main() {
	filename, timerDuration := readArguments()
	file, err := os.Open(filename)

	if err != nil {
		panic("Failed to open " + filename + " file")
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		panic("Failed to read " + filename + " file")
	}

	problems := parseRecords(records)
	correctAnswers := 0

	timer := time.NewTimer(time.Duration(timerDuration) * time.Second)

	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s=", i+1, problem.question)

		answerCh := make(chan string)
		go scanAnswer(answerCh)

		select {
		case <-timer.C:
			fmt.Printf("\nYour time is expired. You scored %d out of %d.", correctAnswers, len(problems))
			return
		case answer := <-answerCh:
			if answer == problem.answer {
				correctAnswers += 1
			}
		}
	}

	fmt.Printf("You scored %d out of %d.", correctAnswers, len(problems))
}

func readArguments() (string, int) {
	filename := flag.String("filename", "problems.csv", "CSV filename")
	timerDuration := flag.Int("time", 30, "Timer duration in seconds")
	flag.Parse()

	return *filename, *timerDuration
}

func parseRecords(records [][]string) []Problem {
	problems := make([]Problem, len(records))

	for i, record := range records {
		problems[i] = Problem{
			question: record[0],
			answer:   record[1],
		}
	}

	return problems
}

func scanAnswer(ch chan string) {
	answer := ""

	if _, err := fmt.Scanln(&answer); err != nil {
		answer = ""
	}

	ch <- answer
}
