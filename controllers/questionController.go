package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

const (
	indexName = "questions_index"
	docType   = "question"
)

type Question struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Distance string `json:"distance"`
}

// Ping function to ping the elasticsearch server
func ping(ctx context.Context, client *elastic.Client, url string) error {
	// Ping the elastic search server
	if client != nil {
		info, code, err := client.Ping(url).Do(ctx)
		if err != nil {
			return err
		}

		fmt.Printf("ElasticSearch returned with code %d and version %s\n", code, info.Version.Number)
		return nil
	}
	return errors.New("Elastic Client is nil")
}

// Create ElasticSearch index if its not exists
func CreateIndexIfDoesNotExists(ctx context.Context, client *elastic.Client, indexName string) error {
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	res, err := client.CreateIndex(indexName).Do(ctx)

	if err != nil {
		return err
	}

	if !res.Acknowledged {
		return errors.New("Created index was not acknowledged. Check that timeout value is correct")
	}
	return nil
}

// Insert a question
func InsertQuestion(ctx context.Context, elasticClient *elastic.Client) {
	// Insert data in elasticSearch
	var questionList []Question
	for index := 1; index < 5; index++ {
		question := Question{
			Title:    fmt.Sprintf("Hey can I use this one, test number %d", index),
			Content:  fmt.Sprintf("U can use it in one condition reje3ha ya weld l3ahira, test number %d", index),
			Distance: fmt.Sprintf("Casablanca, test number %d", index),
		}

		questionList = append(questionList, question)
	}

	for _, questionObj := range questionList {
		_, err := elasticClient.Index().Index(indexName).Type(docType).BodyJson(questionObj).Do(ctx)
		if err != nil {
			fmt.Printf("Question title: %s, Content: %s, Location: %s", questionObj.Title, questionObj.Content, questionObj.Distance)
			continue
		}
	}

	// Flush data (refreshing data in the index) after this we can do get
	elasticClient.Flush().Index(indexName).Do(ctx)
}

// Helper function to convert the response to an Array of questions
func _convertSearchResultToQuestions(searchResult *elastic.SearchResult) []Question {
	var result []Question

	for _, hit := range searchResult.Hits.Hits {
		var questionObj Question
		err := json.Unmarshal(hit.Source, &questionObj)
		if err != nil {
			log.Println("Can't deserialize 'question' object: %s\n", err.Error())
			continue
		}
		result = append(result, questionObj)
	}
	return result
}

// Get all the questions and try to sort them by distance
func GetAll(ctx context.Context, elasticClient *elastic.Client) []Question {
	query := elastic.MatchAllQuery{}

	searchResult, err := elasticClient.Search().Index(indexName).Query(query).Do(ctx)
	if err != nil {
		fmt.Printf("Error during execution GetAll: %s\n", err.Error())
	}

	return _convertSearchResultToQuestions(searchResult)
}
