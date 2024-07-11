package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"news-crawler/internal/models"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSearchClient struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticSearchClient(addresses []string, index string) (*ElasticSearchClient, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ElasticSearchClient{
		client: client,
		index:  index,
	}, nil
}

func (ec *ElasticSearchClient) IndexArticle(article *models.Article) error {
	articleJSON, err := json.Marshal(article)
	if err != nil {
		return err
	}

	_, err = ec.client.Index(
		ec.index,
		bytes.NewReader(articleJSON),
		ec.client.Index.WithContext(context.Background()),
		ec.client.Index.WithDocumentID(article.ID),
	)
	return err
}
