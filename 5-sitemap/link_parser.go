package sitemap

// Solution from Exercise 4.

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// ParseLinks parses the specified HTML file and returns an
// array of Links.
func ParseLinks(r io.Reader) []Link {

	dataHTML, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	var links []Link

	// getLinkNodes traverses through the HTML graph and appends
	// 'link nodes' to `link_nodes`.
	var getLinkNodes func(*html.Node)
	getLinkNodes = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, BuildLink(n))
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getLinkNodes(c)
		}
	}
	getLinkNodes(dataHTML)

	return links
}

// BuildLink returns a Link that is built using the specified
// 'link node'.
func BuildLink(n *html.Node) Link {
	var result Link

	for _, attr := range n.Attr {
		if attr.Key == "href" {
			result.Href = attr.Val
		}
	}

	// getLinkNodeText traverses the sub-nodes of the specified
	// node and returns the text component of the hyperlink.
	var getLinkNodeText func(*html.Node) string
	getLinkNodeText = func(n *html.Node) string {

		if n.Type == html.TextNode {
			return n.Data
		}

		if n.Type != html.ElementNode {
			return ""
		}

		var result string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			result += getLinkNodeText(c) + " "
		}

		return strings.Join(strings.Fields(result), " ")
	}

	result.Text = getLinkNodeText(n)

	return result
}
