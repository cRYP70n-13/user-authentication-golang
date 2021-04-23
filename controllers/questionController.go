package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/teris-io/shortid"
)

const (
	elasticIndexName = "questions"
	elasticTypeName  = "question"
)

type Question struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Distance float64 `josn:"distance"`
}

type QuestionRequest struct {
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Distance float64 `json:"distance"`
}

type QuestionResponse struct {
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Distance float64 `json:"distance"`
}

type SearchResponse struct {
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Hits      string     `json:"hits"`
	Questions []Question `json:"questions"`
}

var elasticClient *elastic.Client

func PostQuestionEndpoint(c *gin.Context) {
	// Parse request
	var qsts []Question

	if err := c.BindJSON(&qsts); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}

	// Inset document in bulk
	bulk := elasticClient.Bulk().Index(elasticIndexName).Type(elasticTypeName)

	for _, q := range qsts {
		qst := Question{
			ID:       shortid.MustGenerate(),
			Title:    q.Title,
			Content:  q.Content,
			Distance: q.Distance,
		}
		bulk.Add(elastic.NewBulkIndexRequest().Id(qst.ID).Doc(qst))
	}
	if _, err := bulk.Do(c.Request.Context()); err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Failed to create questions")
	}
	c.Status(http.StatusOK)
}

func SearchEndpoint(c *gin.Context) {
	// Parse the request
	query := c.Query("query")
	if query == "" {
		errorResponse(c, http.StatusBadRequest, "Query Not specified")
		return
	}

	// Perform the search
	esQuery := elastic.NewMultiMatchQuery(query, "title", "content").Fuzziness("2").MinimumShouldMatch("2")
	result, err := elasticClient.Search().Index(elasticIndexName).Query(esQuery).Do(c.Request.Context())
	if err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	res := SearchResponse{
		Hits: fmt.Sprintf("%d", result.Hits.TotalHits),
	}

	// Transform search result before returning them
	qsts := make([]QuestionResponse, 0)
	for _, hit := range result.Hits.Hits {
		var qst QuestionResponse
		json.Unmarshal(hit.Source, &qst)
		qsts = append(qsts, qst)
	}
	res.Questions = qsts
	c.JSON(http.StatusOK, res)
}

// A helper function to response for errors
func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}
