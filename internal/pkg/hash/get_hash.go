package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetHash(input string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(input))
	if err != nil {
		return "", err
	}
	hashBytes := hash.Sum(nil)
	output := hex.EncodeToString(hashBytes)
	return output, nil
}
