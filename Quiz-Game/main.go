package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "A csv file with 'question, answer' format")
	timeLimt := flag.Int("timeLimit", 30, "The timelimit for the quiz in seconds")
	flag.Parse()
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Error in opening %s file", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse csv file")
	}
	problems := parseLines(lines)
	count := 0
	timer := time.NewTimer(time.Duration(*timeLimt) * time.Second)
	answerChan := make(chan string)
	for i, p := range problems {
		fmt.Printf("Problem %d: %s = ", i+1, p.q)
		go scanAnswer(answerChan)
		select {
		case <-timer.C:
			fmt.Printf("\nYou have scored %d out of %d\n", count, len(problems))
			return
		case answer := <-answerChan:
			if answer == p.a {
				count++
			}
		}
	}
	fmt.Printf("You have scored %d out of %d\n", count, len(problems))
}

func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i].q = strings.TrimSpace(line[0])
		problems[i].a = strings.TrimSpace(line[1])
	}
	return problems
}

func scanAnswer(answerChan chan string) {
	var answer string
	fmt.Scanf("%s", &answer)
	answerChan <- answer
}
