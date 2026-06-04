package inventory

import (
	"inventory-command-module/db/repository"
	domain_inventory "inventory-command-module/internal/domain/inventory"
	"inventory-command-module/internal/domain/shared"
	"time"

	"github.com/iamKienb/go-core/postgres/conv"
	"github.com/jackc/pgx/v5/pgtype"
)

func toDomainInventoryReservations(rows []repository.InventoryReservation) []*domain_inventory.InventoryReservation {
	reservations := make([]*domain_inventory.InventoryReservation, 0, len(rows))
	for _, row := range rows {
		reservations = append(reservations, &domain_inventory.InventoryReservation{
			ID:          row.ID.Bytes,
			InventoryID: row.InventoryID.Bytes,
			SkuID:       row.SkuID.Bytes,
			ShopID:      row.ShopID.Bytes,
			OrderID:     row.OrderID,
			Quantity:    row.Quantity,
			Status:      domain_inventory.ReservationStatus(row.Status),
			ExpiresAt:   row.ExpiresAt.Time,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   &row.UpdatedAt.Time,
		})
	}

	return reservations
}

func toInfraCreateReservationBatch(reservations []*domain_inventory.InventoryReservation) repository.CreateInventoryReservationBatchParams {
	params := repository.CreateInventoryReservationBatchParams{
		Ids:          make([]pgtype.UUID, 0, len(reservations)),
		InventoryIds: make([]pgtype.UUID, 0, len(reservations)),
		OrderIds:     make([]string, 0, len(reservations)),
		Status:       string(reservations[0].Status),
		ExpiresAt:    conv.TimeStampZ(&reservations[0].ExpiresAt),
		CreatedAt:    conv.TimeStampZ(&reservations[0].CreatedAt),
		UpdatedAt:    conv.TimeStampZ(reservations[0].UpdatedAt),
	}

	for _, reservation := range reservations {
		params.Ids = append(params.Ids, conv.UUID(reservation.ID))
		params.InventoryIds = append(params.InventoryIds, conv.UUID(reservation.InventoryID))
		params.OrderIds = append(params.OrderIds, reservation.OrderID)
	}

	return params
}

func toInfraReleaseReservationsByOrderID(params domain_inventory.ReleaseAndFulfillReservationParams) repository.ReleaseReservationsByOrderIDParams {
	return repository.ReleaseReservationsByOrderIDParams{
		OrderID:   params.OrderID,
		UpdatedAt: conv.TimeStampZ(&params.UpdatedAt),
	}
}

func toInfraFulfillReservationsByOrderID(params domain_inventory.ReleaseAndFulfillReservationParams) repository.FulfillReservationsByOrderIDParams {
	return repository.FulfillReservationsByOrderIDParams{
		OrderID:   params.OrderID,
		UpdatedAt: conv.TimeStampZ(&params.UpdatedAt),
	}
}

func toUUUIDs[T ~[16]byte](ids []T) []pgtype.UUID {
	values := make([]pgtype.UUID, 0, len(ids))
	for _, id := range ids {
		values = append(values, conv.UUID(id))
	}

	return values
}

func createdAtOf(inventories []*domain_inventory.Inventory) *time.Time {
	if len(inventories) == 0 {
		return nil
	}

	return &inventories[0].CreatedAt
}

func toUserIDPointer(value pgtype.UUID) *shared.UserID {
	if !value.Valid {
		return nil
	}

	userID := shared.UserID(value.Bytes)
	return &userID
}

func toTimePointer(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}

	timestamp := value.Time
	return &timestamp
}
