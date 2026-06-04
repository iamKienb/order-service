package shared

import "github.com/google/uuid"

var SystemID = NewID[UserID]()

type UserID uuid.UUID

type ShopID uuid.UUID
type ProductID uuid.UUID
type SkuID uuid.UUID
type InventoryID uuid.UUID
type UserAddressID uuid.UUID
type OrderItemID uuid.UUID
type OrderHistoryID uuid.UUID

func NewID[T ~[16]byte]() T {
	return T(uuid.Must(uuid.NewV7()))
}

func (id UserID) String() string {
	return "user_" + uuid.UUID(id).String()
}

func (id ShopID) String() string {
	return "shop_" + uuid.UUID(id).String()
}
func (id ProductID) String() string {
	return "prod_" + uuid.UUID(id).String()
}
func (id SkuID) String() string {
	return "sku_" + uuid.UUID(id).String()
}
func (id InventoryID) String() string {
	return "inv_" + uuid.UUID(id).String()
}
func (id UserAddressID) String() string {
	return "addr_" + uuid.UUID(id).String()
}

func (id OrderItemID) String() string {
	return "it_" + uuid.UUID(id).String()
}

func (id OrderHistoryID) String() string {
	return "history_" + uuid.UUID(id).String()
}

func (id UserID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
func (id ShopID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
func (id ProductID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
func (id SkuID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
func (id InventoryID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
func (id UserAddressID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
func (id OrderItemID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
func (id OrderHistoryID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
