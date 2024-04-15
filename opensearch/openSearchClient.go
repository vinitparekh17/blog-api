package openSearchClient

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/jay-bhogayata/blogapi/logger"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type OpenSearch struct {
	Client *opensearch.Client
}

type OpenSearchQuery struct {
}

type OpenSearchClient interface {
	NewOpenSearchClient(openSearchURL string) (*OpenSearch, error)
	SearchQuery(indexName string, query string) (*opensearchapi.Response, error)
}

type OpenSearchConfig struct {
	URL      string
	UserName string
	Password string
	Env      string
}

func NewOpenSearchClient(osConfig *OpenSearchConfig) (*OpenSearch, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: osConfig.Env == "development"},
		},
		Addresses:     []string{osConfig.URL},
		Username:      osConfig.UserName,
		Password:      osConfig.Password,
		MaxRetries:    3,
		RetryOnStatus: []int{502, 503, 504},
	})

	logger.Log.Info("OpenSearch client is connected")

	return &OpenSearch{Client: client}, err
}

func (o *OpenSearch) SearchQuery(indexName string, query string, ctx context.Context) (*opensearchapi.Response, error) {

	content := strings.NewReader(query)
	search := opensearchapi.SearchRequest{

		Index: []string{indexName},
		Body:  content,
	}

	searchResponse, err := search.Do(ctx, o.Client)
	return searchResponse, err
}
