package grpc_client

import (
	"context"
	"net/http"

	"order-command-module/internal/application/port"
	"order-command-module/internal/domain/shared"

	"connectrpc.com/connect"
	"github.com/iamKienb/api-contract/gen/user"
	"github.com/iamKienb/api-contract/gen/user/userconnect"
)

type userClient struct {
	client userconnect.UserCommandServiceClient
}

func NewUserClient(httpClient *http.Client, baseURL string) port.UserClient {
	return &userClient{client: userconnect.NewUserCommandServiceClient(httpClient, baseURL)}
}

func (c *userClient) GetAddressByID(ctx context.Context, req port.GetAddressByIDRequest) (*port.UserAddress, error) {
	resp, err := c.client.GetUserAddressByID(ctx, connect.NewRequest(&user.GetUserAddressByIDRequest{
		UserId:        req.UserID,
		UserAddressId: req.UserAddressID,
	}))
	if err != nil {
		return nil, err
	}

	address := resp.Msg.GetAddress()
	userID, err := shared.ParseToRawID[shared.UserID](req.UserID)
	if err != nil {
		return nil, err
	}
	addressID, err := shared.ParseToRawID[shared.UserAddressID](address.GetAddressId())
	if err != nil {
		return nil, err
	}

	return &port.UserAddress{
		UserAddressID: addressID,
		UserID:        userID,
		ReceiverName:  address.GetReceiverName(),
		PhoneNumber:   address.GetPhoneNumber(),
		ProvinceID:    int(address.GetProvinceId()),
		ProvinceName:  address.GetProvinceName(),
		WardID:        int(address.GetWardId()),
		WardName:      address.GetWardName(),
		AddressLine:   address.GetAddressLine(),
		Label:         address.GetLabel(),
		IsDefault:     address.GetIsDefault(),
	}, nil
}
