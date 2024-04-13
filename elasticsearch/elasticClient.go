package elasticsearch

import (
	"context"

	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/jay-bhogayata/blogapi/logger"
)

type ElasticClient interface {
	NewElasticClient() ElasticClient
	CreateIndex(index string) error
	IndexDocument(index string, document map[string]interface{}) error
	SearchDocument(index string, query map[string]interface{}) (*search.Response, error)
	BulkIndexDocument(index string, documents []map[string]interface{}) error
}

type elasticClient struct {
	client *elasticsearch8.TypedClient
}

// CreateIndex implements ElasticClient.
func (es *elasticClient) CreateIndex(index string) error {
	_, err := es.client.Indices.Create(index).Do(context.Background())
	return err
}

// IndexDocument implements ElasticClient.
func (es *elasticClient) IndexDocument(index string, document map[string]interface{}) error {
	_, err := es.client.Index(index).
		Request(document).
		Do(context.TODO())

	return err
}

// SearchDocument implements ElasticClient.
func (es *elasticClient) SearchDocument(index string, query map[string]interface{}) (*search.Response, error) {
	result, err := es.client.Search().
		Request(&search.Request{
			Query: &types.Query{
				MatchAll: &types.MatchAllQuery{},
			},
		}).Do(context.TODO())

	return result, err
}

// BulkIndexDocument implements ElasticClient.
func (es *elasticClient) BulkIndexDocument(index string, documents []map[string]interface{}) error {
	multipleDocs := make([]interface{}, len(documents))

	// Converting the documents to interface
	for i, doc := range documents {
		multipleDocs[i] = doc
	}

	_, err := es.client.Bulk().Index(index).Request(&multipleDocs).Do(context.Background())
	return err
}

// NewElasticClient returns a new ElasticClient.
func (*elasticClient) NewElasticClient() ElasticClient {
	es8, err := elasticsearch8.NewTypedClient(elasticsearch8.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		logger.Log.Error("Error occurred while creating new ElasticSearch client", "error", err)
		return nil
	}
	return &elasticClient{client: es8}
}

// NewElasticClient returns a new ElasticClient.
func NewElasticClient() ElasticClient {
	return &elasticClient{}
}
