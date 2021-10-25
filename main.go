package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type WikipediaResponse struct {
	Parse struct {
		Title  string `json:"title"`
		PageId int    `json:"pageid"`
		Text   struct {
			Content string `json:"*"`
		} `json:"text"`
	} `json:"parse"`
}

func main() {
	titleFlag := flag.String("name", "", "name of wikipedia page")
	flag.Parse()
	title := *titleFlag
	escapedTitle := strings.Replace(title, " ", "_", -1)

	response, err := http.Get(fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=parse&section=0&prop=text&format=json&page=%s", escapedTitle))
	if err != nil {
		log.Println(err)
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		log.Println(response.StatusCode)
	}

	var responseBody WikipediaResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		log.Println(err)
		return
	}

	regex := regexp.MustCompile("(?i)" + title)
	matches := regex.FindAllStringIndex(responseBody.Parse.Text.Content, -1)
	fmt.Println(len(matches))
}
