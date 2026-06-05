package service

type Query map[string]any

type QueryBuilder struct {
	must   []Query
	filter []Query
	should []Query
	sorts  []map[string]any
	from   int
	size   int
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		must:   []Query{},
		filter: []Query{},
		should: []Query{},
		sorts:  []map[string]any{},
		size:   defaultPageSize,
	}
}

func (b *QueryBuilder) WithPagination(from, size int) *QueryBuilder {
	if from > 0 {
		b.from = from
	}
	if size > 0 {
		b.size = size
	}
	return b
}

func (b *QueryBuilder) FilterTerm(field string, value any) *QueryBuilder {
	if field == "" || value == nil || value == "" {
		return b
	}
	b.filter = append(b.filter, Query{"term": Query{field: value}})
	return b
}

func (b *QueryBuilder) MustMultiMatch(query string, fields []string) *QueryBuilder {
	if query == "" || len(fields) == 0 {
		return b
	}
	b.must = append(b.must, Query{
		"multi_match": Query{
			"query":  query,
			"fields": fields,
		},
	})
	return b
}

func (b *QueryBuilder) Nested(path string, childBuilder *QueryBuilder) *QueryBuilder {
	if path == "" || childBuilder == nil {
		return b
	}
	childQuery, ok := childBuilder.Build()["query"]
	if !ok {
		return b
	}
	b.must = append(b.must, Query{
		"nested": Query{
			"path":  path,
			"query": childQuery,
		},
	})
	return b
}

func (b *QueryBuilder) WithSort(field, order string) *QueryBuilder {
	if field == "" {
		return b
	}
	if order == "" {
		order = sortAsc
	}
	b.sorts = append(b.sorts, map[string]any{
		field: map[string]any{"order": order},
	})
	return b
}

func (b *QueryBuilder) Build() map[string]any {
	boolQuery := map[string]any{}
	if len(b.must) > 0 {
		boolQuery["must"] = b.must
	}
	if len(b.filter) > 0 {
		boolQuery["filter"] = b.filter
	}
	if len(b.should) > 0 {
		boolQuery["should"] = b.should
	}
	if len(boolQuery) == 0 {
		boolQuery["must"] = []Query{{"match_all": map[string]any{}}}
	}

	body := map[string]any{
		"from":             b.from,
		"size":             b.size,
		"track_total_hits": true,
		"query": map[string]any{
			"bool": boolQuery,
		},
	}
	if len(b.sorts) > 0 {
		body["sort"] = b.sorts
	}
	return body
}
