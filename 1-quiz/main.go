package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

/*
Sources
- https://gosamples.dev/read-csv
- https://www.geeksforgeeks.org/how-to-take-input-from-the-user-in-golang
- https://www.educative.io/answers/how-to-use-the-printf-function-in-golang
- https://www.calhoun.io/6-tips-for-using-strings-in-go/#:~:text=Multiline%20strings,This%20is%20a%20multiline%20string
- https://gobyexample.com/command-line-flags
*/

var DEFAULT_FILEPATH = "problems.csv"

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

// main executes the quiz game.
func main() {

	var filepath string
	flag.StringVar(
		&filepath,
		"filepath",
		DEFAULT_FILEPATH,
		"Filepath (global) to quiz data.",
	)

	flag.Parse()

	quizData := getQuizData(filepath)

	var score int = 0
	var user_response string
	for i, row := range quizData {
		fmt.Printf("%d. %s?\n", i+1, row[0])
		fmt.Scan(&user_response)

		if user_response == row[1] {
			score += 1
		}
	}

	fmt.Printf("\nYou scored %d out of %d.\n", score, len(quizData))
}
