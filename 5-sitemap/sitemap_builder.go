package sitemap

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	URLs []loc  `xml:"url"`
	XML  string `xml:"xml-namespace,attr`
}

// BuildSitemap builds the sitemap for the domain of a
// specified URL.
func BuildSitemap(url string, maxDepth int) urlset {

	urls := GetSitemapURLs(url, maxDepth)

	sitemapXML := urlset{
		XML: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}
	for _, url := range urls {
		sitemapXML.URLs = append(sitemapXML.URLs, loc{url})
	}

	return sitemapXML
}

// GetSitemapURLs traverses the domain of the specified URL
// and returns all pages.
//
// The domain is traversed using the breadth-first search
// algorithm, which is initialized at the base page. This
// implementation is different than the naive BFS algorithm
// to account for cycling.
func GetSitemapURLs(url string, maxDepth int) []string {

	// Track which URLs have been visited to prevent cycling.
	visited := map[string]bool{}

	queue := map[string]bool{}
	nextQueue := map[string]bool{url: true}

	for i := 0; i < maxDepth; i++ {
		queue, nextQueue = nextQueue, map[string]bool{}

		for url, _ := range queue {
			if _, ok := visited[url]; ok {
				continue
			}

			// Flag URL.
			visited[url] = true

			// Enqueue child URLs.
			for _, link := range GetLinksFromPage(url) {
				nextQueue[link] = true
			}

		}
	}

	// Create an array of URLs that were visited during
	// execution.
	var urls []string
	for url, _ := range visited {
		urls = append(urls, url)
	}

	return urls
}

// GetLinksFromPage retrieves a web page via HTTP and returns
// all links within it's domain.
func GetLinksFromPage(page string) []string {

	resp, err := http.Get(page)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	links := ParseLinks(resp.Body)

	var hrefs []string
	for _, l := range links {

		href := l.Href
		if strings.HasPrefix(href, "/") {
			href = base + l.Href
		}

		if strings.HasPrefix(href, base) {
			hrefs = append(hrefs, href)
		}
	}

	return hrefs
}
