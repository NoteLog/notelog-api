package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

func GetGitHubRepo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	index := "repo*"
	field := "*"

	res := esSearchGitHub(index, field, query)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func esSearchGitHub(index, field, term string) []byte {
	esURL := os.Getenv("ESURL")
	esPWD := os.Getenv("ESPWD")

	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
		Username:  "elastic",
		Password:  esPWD,
	}
	es, err := elasticsearch.NewClient(cfg)

	query := fmt.Sprintf(`{
		  "query": {
		    "multi_match": {
		      "query": "%s",
		      "type": "most_fields",
		      "fields": [
		        "username",
		        "reponame",
				"description"
		      ],
		      "operator": "and",
			  "fuzziness": "AUTO"
		    }
		  },
		"highlight": {
			"fields": {
				"username": {},
				"reponame": {},
				"description": {}
			},
			"pre_tags": ["<b>"],
			"post_tags": ["</b>"]
		  }
		}`, term)
	queryString := strings.NewReader(query)

	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(queryString),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	return resBody
}
