package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
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
var DEFAULT_SHUFFLE = false

// main executes the quiz game.
func main() {

	var filepath string
	flag.StringVar(
		&filepath,
		"filepath",
		DEFAULT_FILEPATH,
		"Filepath (global) to quiz data.",
	)

	var shuffleFlag bool
	flag.BoolVar(
		&shuffleFlag,
		"shuffle",
		DEFAULT_SHUFFLE,
		"Shuffle the quiz data.",
	)

	var timeLimit int
	flag.IntVar(
		&timeLimit,
		"time_limit",
		DEFAULT_TIME_LIMIT,
		"Quiz duration (in seconds).",
	)

	flag.Parse()

	data := getCSVData(filepath)

	shuffleText := "Shuffling is off.\n"
	if shuffleFlag {
		data = shuffleCSVData(data)
		shuffleText = "Shuffling is on.\n"
	}

	fmt.Printf(
		strings.Join(
			[]string{
				"You are about to take a quiz.\n",
				"The timer is set to %d seconds.\n",
				shuffleText,
				"Please press enter to begin.\n",
			},
			"",
		),
		timeLimit,
	)

	userResponse := getUserInput()

	if userResponse != "\n" {
		fmt.Println("Text detected; terminating quiz.")
		os.Exit(0)
	}

	score := executeQuiz(data, timeLimit)

	fmt.Printf("\nYou scored %d out of %d.\n", score, len(data))
}

// getUserInput uses Buffered I/O to read in user input from the command line.
// This allows for user input with spaces.
func getUserInput() string {

	in := bufio.NewReader(os.Stdin)

	result, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return result
}

// getCSVData reads in a CSV file from the specified filepath
// and returns the data.
func getCSVData(filepath string) [][]string {

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

// shuffleCSVData returns a copy of the data but shuffled
// psuedo-randomly using the default seed.
func shuffleCSVData(CSVData [][]string) [][]string {

	var data [][]string
	copy(data, CSVData)

	rand.Shuffle(
		len(data),
		func(i int, j int) {
			data[i], data[j] = data[j], data[i]
		},
	)

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

	resultChannel := make(chan int, 1)
	timeoutChannel := time.After(time.Duration(timeLimit) * time.Second)

	score := 0
	for i, row := range data {

		fmt.Printf("%d. %s?\n", i+1, row[0])

		// A goroutine that prompts the user with a quiz question.
		//
		// The user response and the answer to the corresponding quiz question
		// are stripped of formatting to ensure a valid comparison.
		//
		// This is done in a goroutine to allow for a timeout.
		go func() {

			userResponse := getUserInput()

			processedUserResponse := strings.ReplaceAll(
				strings.Trim(strings.ToLower(userResponse), "\n"), " ", "",
			)

			processedAnswer := strings.ReplaceAll(
				strings.ToLower(row[1]), " ", "",
			)

			result := 0
			if processedUserResponse == processedAnswer {
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
