package database

import (
	"fmt"

	"github.com/olivere/elastic/v7"
)

func GetESClient() (*elastic.Client, error) {
	Client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)

	fmt.Println("ElasticSearch initialized ..")
	return Client, err
}
