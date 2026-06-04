package shared

import (
	"crypto/rand"
	"math/big"
)

const base36Charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CryptoRandomString(length int) string {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(base36Charset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			result[i] = base36Charset[0]
			continue
		}
		result[i] = base36Charset[num.Int64()]
	}

	return string(result)
}
