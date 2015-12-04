package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"
)

const (
	rubyGemsUrl = "https://rubygems.org/api/v1/search.json?query=%s"
	noGemFound  = "No gem found"
	oops        = "Oops..."
	icon        = "icon.png"
	cmdMod      = "cmd"
)

type gem struct {
	XMLName    xml.Name   `xml:"item"`
	Name       string     `xml:"title"`
	ProjectUri string     `xml:"arg,attr"`
	Subtitles  []subtitle `xml:"subtitle"`
	Icon       string
}

type subtitle struct {
	Mod     string `xml:"mod,attr,omitempty"`
	Content string `xml:",chardata"`
}

type gemResults struct {
	XMLName xml.Name `xml:"items"`
	Gems    []gem
}

func main() {
	searchString := flag.String("search", "", "This is the gem your searching for")
	flag.Parse()

	url := fmt.Sprintf(rubyGemsUrl, *searchString)

	resp, err_conn := http.Get(url)
	if err_conn != nil {
		outputError(url)
		return
	}

	dec := json.NewDecoder(resp.Body)
	var json []interface{}
	err_decode := dec.Decode(&json)
	if err_decode != nil {
		outputError(url)
		return
	}

	gems := make([]gem, len(json))
	for i := 0; i < len(json); i++ {
		gems[i] = convertJsonToGem(json[i].(map[string]interface{}))
	}

	if len(gems) == 0 {
		outputError(url)
		return
	}

	outputXML(gems)
}

func outputError(url string) {
	outputXML(createErrorGemResponse(url))
}

func createErrorGemResponse(url string) []gem {
	return []gem{gem{Name: oops, Subtitles: []subtitle{{Content: noGemFound}}, ProjectUri: url, Icon: icon}}
}

func outputXML(gems []gem) {
	gemResults := gemResults{Gems: gems}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("  ", "    ")
	if err := enc.Encode(gemResults); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func convertJsonToGem(input map[string]interface{}) gem {
	subttitles := []subtitle{
		{Content: fmt.Sprintf("%s, Version: %s", input["info"].(string), input["version"].(string))},
		{Mod: cmdMod, Content: input["homepage_uri"].(string)},
	}

	return gem{
		Name:       input["name"].(string),
		ProjectUri: input["project_uri"].(string),
		Subtitles:  subttitles,
		Icon:       icon,
	}
}
