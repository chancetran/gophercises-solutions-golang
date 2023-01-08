package main

import (
	"cyoa"
	"fmt"
	"net/http"
)

// main executes a local server that hosts an interactive
// 'choose your own adventure' story via HTML.
func main() {

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(
		":8080",
		cyoa.AdventureHandler{
			cyoa.GetAdventureFromJSONFile("gopher.json"),
		},
	)

}
