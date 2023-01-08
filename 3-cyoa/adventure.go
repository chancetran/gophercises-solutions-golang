package cyoa

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Adventure map[string]Event

type Event struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

// GetAdventureFromJSONFile reads in JSON data from the
// specified filepath.
//
// The JSON data is unmarshalled into an Adventure, or a map
// of story-keys to Events. These events contain an array of
// options which drive the adventure.
func GetAdventureFromJSONFile(filepath string) Adventure {

	jsn, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	adventure := Adventure{}
	err = json.Unmarshal(jsn, &adventure)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return adventure
}
