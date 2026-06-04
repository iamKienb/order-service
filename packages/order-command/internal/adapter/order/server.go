package inventory

import (
	"context"

	"inventory-command-module/internal/application/commands/create_inventories"
	"inventory-command-module/internal/application/commands/delete_inventories"
	"inventory-command-module/internal/application/commands/fulfill_stock"
	"inventory-command-module/internal/application/commands/release_stock"
	"inventory-command-module/internal/application/commands/reserve_stock"

	"connectrpc.com/connect"
	"github.com/iamKienb/api-contract/gen/inventory"
	"github.com/iamKienb/api-contract/gen/inventory/inventoryconnect"
	"github.com/iamKienb/go-core/app_error"
	authx "github.com/iamKienb/go-core/middleware/auth"
)

type inventoryServer struct {
	createInventoriesExecutor create_inventories.Executor
	deleteInventoriesExecutor delete_inventories.Executor
	reserveStockExecutor      reserve_stock.Executor
	releaseStockExecutor      release_stock.Executor
	fulfillStockExecutor      fulfill_stock.Executor
}

func NewInventoryServer(
	createInventoriesExecutor create_inventories.Executor,
	deleteInventoriesExecutor delete_inventories.Executor,
	reserveStockExecutor reserve_stock.Executor,
	releaseStockExecutor release_stock.Executor,
	fulfillStockExecutor fulfill_stock.Executor,
) *inventoryServer {
	return &inventoryServer{
		createInventoriesExecutor: createInventoriesExecutor,
		deleteInventoriesExecutor: deleteInventoriesExecutor,
		reserveStockExecutor:      reserveStockExecutor,
		releaseStockExecutor:      releaseStockExecutor,
		fulfillStockExecutor:      fulfillStockExecutor,
	}
}

func (s *inventoryServer) CreateInventories(ctx context.Context, req *connect.Request[inventory.CreateInventoriesRequest]) (*connect.Response[inventory.CreateInventoriesResponse], error) {
	currentUser, err := requireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := ToCreateInventoriesCommand(currentUser.UserID, req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.createInventoriesExecutor.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(ToCreateInventoriesResponse(result)), nil
}

func (s *inventoryServer) DeleteInventory(ctx context.Context, req *connect.Request[inventory.DeleteInventoryRequest]) (*connect.Response[inventory.DeleteInventoryResponse], error) {
	currentUser, err := requireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := ToDeleteInventoriesCommand(currentUser.UserID, req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.deleteInventoriesExecutor.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(ToDeleteInventoriesResponse(result)), nil
}

func (s *inventoryServer) ReserveStock(ctx context.Context, req *connect.Request[inventory.ReserveStockRequest]) (*connect.Response[inventory.ReserveStockResponse], error) {
	currentUser, err := requireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := ToReserveStockCommand(currentUser.UserID, req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.reserveStockExecutor.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(ToReserveStockResponse(result)), nil
}

func (s *inventoryServer) ReleaseStock(ctx context.Context, req *connect.Request[inventory.ReleaseStockRequest]) (*connect.Response[inventory.ReleaseStockResponse], error) {
	currentUser, err := requireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := ToReleaseStockCommand(currentUser.UserID, req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.releaseStockExecutor.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(ToReleaseStockResponse(result)), nil
}

func (s *inventoryServer) FulfillStock(ctx context.Context, req *connect.Request[inventory.FulfillStockRequest]) (*connect.Response[inventory.FulfillStockResponse], error) {
	currentUser, err := requireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := ToFulfillStockCommand(currentUser.UserID, req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.fulfillStockExecutor.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(ToAFulfilStockResponse(result)), nil
}

func requireCurrentUser(ctx context.Context) (*authClaims, error) {
	claims := authx.GetUserInfoFromCtx(ctx)
	if claims == nil || claims.UserID == "" {
		return nil, app_error.Unauthorized("authentication required")
	}

	return &authClaims{UserID: claims.UserID}, nil
}

type authClaims struct {
	UserID string
}

var _ inventoryconnect.InventoryCommandServiceHandler = (*inventoryServer)(nil)
