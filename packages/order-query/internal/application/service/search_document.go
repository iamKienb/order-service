package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type SearchHit[T any] struct {
	Source T
}

type SearchResult[T any] struct {
	Hits  []SearchHit[T]
	Total int64
}

func SearchDocuments[T any](ctx context.Context, esClient *elasticsearch.TypedClient, index string, queryBody map[string]any) (*SearchResult[T], error) {
	rawJSON, err := json.Marshal(queryBody)
	if err != nil {
		return nil, fmt.Errorf("marshal elasticsearch query: %w", err)
	}

	response, err := esClient.Search().Index(index).Raw(bytes.NewReader(rawJSON)).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("search elasticsearch index %s: %w", index, err)
	}

	results := make([]SearchHit[T], 0, len(response.Hits.Hits))
	for _, hit := range response.Hits.Hits {
		var source T
		if err := json.Unmarshal(hit.Source_, &source); err != nil {
			return nil, fmt.Errorf("decode elasticsearch source: %w", err)
		}
		results = append(results, SearchHit[T]{Source: source})
	}

	return &SearchResult[T]{
		Hits:  results,
		Total: totalHits(response.Hits),
	}, nil
}

func totalHits(hits types.HitsMetadata) int64 {
	if hits.Total == nil {
		return 0
	}
	return hits.Total.Value
}
