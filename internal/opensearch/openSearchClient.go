package openSearchClient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jay-bhogayata/blogapi/internal/logger"
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
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: osConfig.Env == "development",
				Rand:               nil,
				MinVersion:         tls.VersionTLS12,
			},
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

func (o *OpenSearch) QueryBuilder(matchType string, fields []string, term string) (string, error) {
	if len(fields) == 0 {
		return "", fmt.Errorf("fields cannot be empty")
	}

	var queryBody map[string]interface{}

	if matchType == "multi" {
		// For "multi" match
		queryBody = map[string]interface{}{
			"query": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  term,
					"fields": fields,
				},
			},
		}
	} else {
		// For "single" match (default)
		queryBody = map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					fields[0]: term,
				},
			},
		}
	}

	// Convert map to JSON
	queryJSON, err := json.Marshal(queryBody)
	if err != nil {
		return "", err
	}

	return string(queryJSON), nil
}
