package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"user-athentication-golang/models"
	routes "user-athentication-golang/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/olivere/elastic/v7"
	"github.com/teris-io/shortid"
)

const (
	elasticIndexName = "question_index"
	elasticIndexType = "question"
)

var (
	elasticClient *elastic.Client
)

// Helper function to make costume errors
func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}

// The function to create a questions
func CreateQuestionEndpoint(c *gin.Context) {
	var questions []models.QuestionRequest
	if err := c.BindJSON(&questions); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}

	bulk := elasticClient.Bulk().Index(elasticIndexName).Type(elasticIndexType)
	for _, d := range questions {
		qst := models.Question{
			ID:        shortid.MustGenerate(),
			Title:     d.Title,
			CreatedAt: time.Now().UTC(),
			Content:   d.Content,
		}
		bulk.Add(elastic.NewBulkIndexRequest().Id(qst.ID).Doc(qst))
	}
	if _, err := bulk.Do(c.Request.Context()); err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Failed to post questions")
		return
	}
	c.Status(http.StatusOK)
}

// The function to create a questions
func PostAnswer(c *gin.Context) {
	var answers []models.Answer
	if err := c.BindJSON(&answers); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}

	bulk := elasticClient.Bulk().Index(elasticIndexName).Type(elasticIndexType)
	for _, d := range answers {
		ans := models.Answer{
			ID:     shortid.MustGenerate(),
			Title:  d.Title,
			Answer: d.Answer,
		}
		bulk.Add(elastic.NewBulkIndexRequest().Id(ans.ID).Doc(ans))
	}
	if _, err := bulk.Do(c.Request.Context()); err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Failed to post answers")
		return
	}
	c.Status(http.StatusOK)
}

// The search Endpoint function
func searchEndpoint(c *gin.Context) {
	// Parse the request
	query := c.Query("query")
	if query == "" {
		// errorResponse(c, http.StatusBadRequest, "Query not specified")
		query = `{
			"query": {
			  "match_all": {}
			}
		  }`
	}

	skip := 0
	take := 10
	if i, err := strconv.Atoi(c.Query("skip")); err == nil {
		skip = i
	}
	if i, err := strconv.Atoi(c.Query("take")); err == nil {
		take = i
	}

	esQuery := elastic.NewMultiMatchQuery(query, "title", "content").
		Fuzziness("2").
		MinimumShouldMatch("2")
	result, err := elasticClient.Search().
		Index(elasticIndexName).
		Query(esQuery).
		From(skip).Size(take).
		Do(c.Request.Context())

	if err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	res := models.SearchResponse{
		Time: fmt.Sprintf("%d", result.TookInMillis),
		Hits: fmt.Sprintf("%d", result.Hits.TotalHits),
	}

	// The questions
	docs := make([]models.QuestionResponse, 0)
	for _, hit := range result.Hits.Hits {
		var doc models.QuestionResponse
		json.Unmarshal(hit.Source, &doc)
		docs = append(docs, doc)
	}
	res.Questions = docs
	c.JSON(http.StatusOK, res)
}

func main() {
	port := os.Getenv("PORT")
	var err error
	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL("http://localhost:9200"),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		}
		break
	}

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)

	// Elastic Search part
	router.POST("/api/v1/questions", CreateQuestionEndpoint)
	router.POST("/api/v1/answers", PostAnswer)
	router.GET("/api/v1/search", searchEndpoint)
	if err = router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	router.Run(":" + port)
}
