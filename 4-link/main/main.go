package main

import (
	"flag"
	"fmt"
	"link"
)

/*
Source
- https://pkg.go.dev/golang.org/x/net/html
*/

// main parses the HTML file specified with the flag
// `html_path`, or the default file, and prints all links
//  with their corresponding texts.
func main() {

	var html_path string
	flag.StringVar(
		&html_path, "html_path", "data/ex1.html", "path to HTML file",
	)

	flag.Parse()

	p := link.ParseLinksFromHTTP(html_path)

	for _, link := range p {
		fmt.Printf(
			"{href: %s\n text: %s}\n", link.Href, link.Text,
		)
	}

}
