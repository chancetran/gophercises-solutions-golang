package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"sitemap"
)

/*
Source
- https://courses.calhoun.io/lessons/les_goph_24
- https://www.sitemaps.org/
- https://en.wikipedia.org/wiki/Breadth-first_search
- https://pkg.go.dev/encoding/xml#Encoder.Encode
*/

// main builds and prints a sitemap in XML for the domain of
// a specified URL.
func main() {

	var baseURL string
	flag.StringVar(
		&baseURL,
		"url",
		"https://gophercises.com",
		"URL to build a sitemap for",
	)

	var maxDepth int
	flag.IntVar(
		&maxDepth,
		"max_depth",
		3,
		"maximum number of subsequent links to traverse",
	)

	flag.Parse()

	sitemapXML := sitemap.BuildSitemap(baseURL, maxDepth)

	// Print the sitemap in the terminal.
	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(sitemapXML); err != nil {
		panic(err)
	}
	fmt.Println()
}
