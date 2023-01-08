package main

import (
	"cyoa"
	"fmt"
	"net/http"
)

/*
Sources
- https://medium.com/geekculture/demystifying-http-handlers-in-golang-a363e4222756
- https://stackoverflow.com/questions/34031801/function-declaration-syntax-things-in-parenthesis-before-function-name
- https://pkg.go.dev/html/template
- https://gowebexamples.com/templates/
- https://www.w3schools.com/html/html_intro.asp
- https://www.w3schools.com/tags/att_a_href.asp
- https://stackoverflow.com/questions/2906582/how-do-i-create-an-html-button-that-acts-like-a-link
*/

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
