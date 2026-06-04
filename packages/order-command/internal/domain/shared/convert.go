package shared

import (
	"strings"

	"github.com/google/uuid"
)

func ParseToRawID[T ~[16]byte](id string) (T, error) {
	var result T
	parts := strings.Split(id, "_")
	rawID := parts[len(parts)-1]

	parsed, err := uuid.Parse(rawID)
	if err != nil {
		return result, err
	}

	return T(parsed), nil
}
