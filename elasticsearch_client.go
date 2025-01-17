package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

const (
	ES_URL = "http://IP_ADDRES:9200"
)

func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("elastic_username", "elasticsearch_password"))
	if err != nil {
		panic(err)
	}

	searchResult, err := client.Search().
		Index(index).
		Query(query).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func saveToES(i interface{}, index string, id string) error {
	client, err := elastic.NewClient(
		elastic.SetURL(ES_URL),
		elastic.SetBasicAuth("elastic_username", "elasticsearch_password"))
	if err != nil {
		return err
	}

	_, err = client.Index().Index(index).Id(id).BodyJson(i).Do(context.Background())

	if err != nil {
		return err
	}

	fmt.Printf("Data is saved to Elasticsearch: %s\n", id)
	return nil
}
