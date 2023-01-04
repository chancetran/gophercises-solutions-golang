package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

/*
Sources
- https://gosamples.dev/read-csv
- https://www.geeksforgeeks.org/how-to-take-input-from-the-user-in-golang
- https://www.educative.io/answers/how-to-use-the-printf-function-in-golang
- https://www.calhoun.io/6-tips-for-using-strings-in-go/#:~:text=Multiline%20strings,This%20is%20a%20multiline%20string
- https://gobyexample.com/command-line-flags
- https://go.dev/tour/concurrency/1
- https://gobyexample.com/timeouts
*/

var DEFAULT_FILEPATH = "problems.csv"
var DEFAULT_TIME_LIMIT = 30

// getQuizData reads in a file from the specified filepath and returns the
// data.
func getQuizData(filepath string) [][]string {

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return data
}

// executeQuiz loops over the specified data, displaying the first
// row-attribute as the 'question', awaits user input, and then compares said
// input with the second row-attribute as the 'answer'.
//
// If the user input matches the 'answer' the running total, or score, of
// correct answers is incremented. Otherwise, the loop continues.
//
// The loop is terminated if the function execution time exceeds the
// 'timeLimit'.
func executeQuiz(data [][]string, timeLimit int) int {

	var userResponse string

	resultChannel := make(chan int, 1)
	timeoutChannel := time.After(time.Duration(timeLimit) * time.Second)

	score := 0
	for i, row := range data {
		fmt.Printf("%d. %s?\n", i+1, row[0])

		go func() {
			// TODO: CHANGE TO BUFIO SO IT CAN ACCEPT STRINGS WITH SPACES.
			fmt.Scan(&userResponse)

			result := 0
			if userResponse == row[1] {
				result = 1
			}

			resultChannel <- result
		}()

		select {
		case result := <-resultChannel:
			score += result
		case <-timeoutChannel:
			fmt.Println("Time's up!")
			return score
		}

	}

	return score
}

// main executes the quiz game.
func main() {

	var filepath string
	flag.StringVar(
		&filepath,
		"filepath",
		DEFAULT_FILEPATH,
		"Filepath (global) to quiz data.",
	)

	var timeLimit int
	flag.IntVar(
		&timeLimit,
		"time_limit",
		DEFAULT_TIME_LIMIT,
		"Quiz duration (in seconds).",
	)

	flag.Parse()

	data := getQuizData(filepath)

	fmt.Printf(
		strings.Join(
			[]string{
				"You are about to take a quiz.\n",
				"The timer is set to %d seconds.\n",
				"Please press enter to begin.\n",
			},
			"",
		),
		timeLimit,
	)

	var userResponse string

	// TODO: CHANGE TO BUFIO SO IT CAN ACCEPT STRINGS WITH SPACES.
	fmt.Scanln(&userResponse)

	if userResponse != "" {
		fmt.Println("Text detected; terminating quiz.")
		os.Exit(0)
	}

	score := executeQuiz(data, timeLimit)

	fmt.Printf("\nYou scored %d out of %d.\n", score, len(data))
}
