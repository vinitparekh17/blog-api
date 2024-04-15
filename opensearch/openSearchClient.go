package openSearchClient

import (
	"context"
	"crypto/tls"
	"fmt"
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

	res, pingErr := client.Ping()
	if err != nil {
		return nil, pingErr
	}
	defer res.Body.Close()

	logger.Log.Info(fmt.Sprintf("OpenSearch client is connected to %s\nPing response: %s", osConfig.URL, res.Status()))

	return &OpenSearch{Client: client}, err
}

func (o *OpenSearch) SearchQuery(indexName string, query string) (*opensearchapi.Response, error) {
	search := opensearchapi.SearchRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(query),
	}

	searchResponse, err := search.Do(context.TODO(), o.Client)
	return searchResponse, err
}
