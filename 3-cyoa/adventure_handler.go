package cyoa

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type AdventureHandler struct {
	Adventure Adventure
}

// ServeHTTP for AdventureHandler Handlers fetches the current
// Event specified in the request, or the introduction Event
// if one is not, and generates an HTML page using a template.
func (ah AdventureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestedPath := r.URL.String()

	if requestedPath == "/" {
		requestedPath = "/intro"
	}
	requestedPath = strings.TrimPrefix(requestedPath, "/")

	if event, ok := ah.Adventure[requestedPath]; ok {
		GetEventPage(w, r, event)
	}
}

// GetEventPage generates an HTML page for the given event
// using an HTML template.
func GetEventPage(w http.ResponseWriter, r *http.Request, data Event) error {

	const EVENT_TEMPLATE = `
	<h1>{{.Title}}</h1>
	<body>
		{{range .Story}}
			<p>
				{{.}}
			</p>
		{{end}}
		{{range .Options}}
			<a href="{{.Arc}}" class="button">
				{{.Text}}
			</a> 
			<br>
		{{end}}
	</body>
	`

	tmpl, err := template.New("EventTemplate").Parse(EVENT_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, data)

	return err
}
