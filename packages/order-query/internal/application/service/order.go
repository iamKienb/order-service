package service

import (
	"context"
	"strconv"

	"order-query-module/internal/application/queries/get_order_detail"
	"order-query-module/internal/application/queries/list_buyer_orders"
	"order-query-module/internal/application/queries/list_shop_orders"
	"order-query-module/internal/application/queries/search_orders"
	"order-query-module/internal/application/service/models"
	"order-shared-module/alias"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/iamKienb/go-core/app_error"
	esx "github.com/iamKienb/go-core/elasticsearch"
)

type QueryService interface {
	GetOrderDetail(context.Context, get_order_detail.Query) (*get_order_detail.Result, error)
	ListBuyerOrders(context.Context, list_buyer_orders.Query) (*list_buyer_orders.Result, error)
	ListShopOrders(context.Context, list_shop_orders.Query) (*list_shop_orders.Result, error)
	SearchOrders(context.Context, search_orders.Query) (*search_orders.Result, error)
}

type queryService struct {
	esClient *elasticsearch.TypedClient
	index    string
}

const (
	defaultPageSize = 20
	maxPageSize     = 100
	sortAsc         = "asc"
	sortDesc        = "desc"
	errOrderMissing = "order was not found"
)

func NewQueryService(esService esx.ESXService) QueryService {
	return &queryService{
		esClient: esService.GetClient(),
		index:    alias.OrderAlias,
	}
}

func (s *queryService) GetOrderDetail(ctx context.Context, query get_order_detail.Query) (*get_order_detail.Result, error) {
	searchQuery := NewQueryBuilder().
		WithPagination(0, 1).
		FilterTerm("id", query.OrderID).
		Build()

	result, err := SearchDocuments[models.Order](ctx, s.esClient, s.index, searchQuery)
	if err != nil {
		return nil, err
	}
	if len(result.Hits) == 0 {
		return nil, app_error.NotFound(errOrderMissing, nil)
	}

	return &get_order_detail.Result{Order: &result.Hits[0].Source}, nil
}

func (s *queryService) ListBuyerOrders(ctx context.Context, query list_buyer_orders.Query) (*list_buyer_orders.Result, error) {
	page := normalizePage(query.Page)
	builder := NewQueryBuilder().
		WithPagination(pageOffset(page.Token), page.Size).
		FilterTerm("buyer_id", query.BuyerID).
		FilterTerm("status", query.Status).
		WithSort("created_at", sortDesc).
		WithSort("id", sortAsc)

	return s.searchPage(ctx, builder.Build(), page)
}

func (s *queryService) ListShopOrders(ctx context.Context, query list_shop_orders.Query) (*list_shop_orders.Result, error) {
	page := normalizePage(query.Page)
	builder := NewQueryBuilder().
		WithPagination(pageOffset(page.Token), page.Size).
		FilterTerm("shop_id", query.ShopID).
		FilterTerm("status", query.Status).
		WithSort("created_at", sortDesc).
		WithSort("id", sortAsc)

	return s.searchPage(ctx, builder.Build(), page)
}

func (s *queryService) SearchOrders(ctx context.Context, query search_orders.Query) (*search_orders.Result, error) {
	page := normalizePage(query.Page)
	builder := NewQueryBuilder().
		WithPagination(pageOffset(page.Token), page.Size).
		FilterTerm("shop_id", query.ShopID).
		FilterTerm("buyer_id", query.BuyerID).
		FilterTerm("status", query.Status).
		WithSort("created_at", sortDesc).
		WithSort("order_id", sortAsc)

	if query.Keyword != "" {
		itemQuery := NewQueryBuilder().MustMultiMatch(query.Keyword, []string{"items.product_name"})
		builder.Nested("items", itemQuery)
	}

	return s.searchPage(ctx, builder.Build(), page)
}

func (s *queryService) searchPage(ctx context.Context, queryBody map[string]any, page models.Page) (*models.OrderPage, error) {
	result, err := SearchDocuments[models.Order](ctx, s.esClient, s.index, queryBody)
	if err != nil {
		return nil, err
	}

	items := ordersFromHits(result.Hits)
	return &models.OrderPage{
		Items:         items,
		Total:         result.Total,
		NextPageToken: nextPageToken(page, len(items), result.Total),
	}, nil
}

func ordersFromHits(hits []SearchHit[models.Order]) []models.Order {
	items := make([]models.Order, 0, len(hits))
	for _, hit := range hits {
		items = append(items, hit.Source)
	}
	return items
}

func normalizePage(page models.Page) models.Page {
	if page.Size <= 0 || page.Size > maxPageSize {
		page.Size = defaultPageSize
	}
	return page
}

func pageOffset(token string) int {
	if token == "" {
		return 0
	}
	offset, err := strconv.Atoi(token)
	if err != nil || offset < 0 {
		return 0
	}
	return offset
}

func nextPageToken(page models.Page, resultCount int, total int64) string {
	nextOffset := pageOffset(page.Token) + resultCount
	if int64(nextOffset) >= total || resultCount == 0 {
		return ""
	}
	return strconv.Itoa(nextOffset)
}
