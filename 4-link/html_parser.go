package link

import (
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// ParseLinksFromHTTP parses the specified HTML file and
// returns an array of Links.
func ParseLinksFromHTTP(filepath string) []Link {

	var links []Link

	data, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	dataHTML, err := html.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	defer data.Close()

	var link_nodes []*html.Node

	// getLinkNodes traverses through the HTML graph and appends
	// 'link nodes' to `link_nodes`.
	var getLinkNodes func(*html.Node)
	getLinkNodes = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link_nodes = append(link_nodes, n)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getLinkNodes(c)
		}
	}
	getLinkNodes(dataHTML)

	for _, link_node := range link_nodes {
		links = append(links, BuildLink(link_node))
	}

	return links
}

// BuildLink returns a Link using the specified 'link node'.
func BuildLink(n *html.Node) Link {
	var result Link

	for _, attr := range n.Attr {
		if attr.Key == "href" {
			result.Href = attr.Val
		}
	}

	// getLinkNodeText traverses the sub-nodes of the specified
	// node and returns the text component of a hyperlink.
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
