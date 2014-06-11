package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type gem struct {
	XMLName    xml.Name `xml:"item"`
	Name       string   `xml:"title"`
	ProjectUri string   `xml:"arg,attr"`
	Info       string   `xml:"subtitle"`
	Icon       string
}
type gemResults struct {
	XMLName xml.Name `xml:"items"`
	Gems    []gem
}

func main() {
	searchString := flag.String("search", "", "This is the gem your searching for")
	flag.Parse()

	url := fmt.Sprintf("http://rubygems.org/api/v1/search.json?query=%s", *searchString)

	resp, err_conn := http.Get(url)
	if err_conn != nil {
		outputError(url)
		//fmt.Println("Something happened when trying to query ruby gems: ", err_conn)
		return
	}

	dec := json.NewDecoder(resp.Body)
	var json []interface{}
	err_decode := dec.Decode(&json)
	if err_decode != nil {
		outputError(url)
		//fmt.Println("Error occured when trying to decode json body: ", err_decode)
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
	return []gem{gem{Name: "Oops...", Info: "No gem found", ProjectUri: url, Icon: "icon.png"}}
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
	return gem{
		Name:       input["name"].(string),
		ProjectUri: input["project_uri"].(string),
		Info:       fmt.Sprintf("%s, Version: %s", input["info"].(string), input["version"].(string)),
		Icon:       "icon.png",
	}
}
