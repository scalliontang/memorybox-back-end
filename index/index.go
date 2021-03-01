package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

const (
	POST_INDEX = "post"
	USER_INDEX = "user"
	ES_URL     = "http://IP_ADDRES:9200"
)

func main() {
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("elastic", "your_elasticsearch_password"))
	if err != nil {
		panic(err)
	}
	exists, err := client.IndexExists(POST_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		mapping := `{
                        "mappings": {
                                "properties": {
                                        "user":     { "type": "keyword" },
                                        "message":  { "type": "text" },
                                        "url":      { "type": "keyword", "index": false },
                                        "type":     { "type": "keyword", "index": false }
                                }
                        }
                }`
		_, err := client.CreateIndex(POST_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	exists, err = client.IndexExists(USER_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		mapping := `{
                        "mappings": {
                                "properties": {
                                        "username": {"type": "keyword"},
                                        "password": {"type": "keyword", "index": false},
                                        "age":      {"type": "long", "index": false},
                                        "gender":   {"type": "keyword", "index": false}
                                }
                        }
                }`
		_, err = client.CreateIndex(USER_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Indexes are created.")
}
