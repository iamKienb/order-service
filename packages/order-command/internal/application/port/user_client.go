package port

import (
	"context"
	"order-command-module/internal/domain/shared"
)

type GetAddressByIDRequest struct {
	UserID        string
	UserAddressID string
}

type UserAddress struct {
	UserAddressID shared.UserAddressID
	UserID        shared.UserID
	ReceiverName  string
	PhoneNumber   string
	ProvinceID    int
	ProvinceName  string
	WardID        int
	WardName      string
	AddressLine   string
	Label         string
	IsDefault     bool
}

type UserClient interface {
	GetAddressByID(ctx context.Context, req GetAddressByIDRequest) (*UserAddress, error)
}
