package order

import "github.com/google/uuid"

func orderAggregateID(orderID string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(orderID))
}
