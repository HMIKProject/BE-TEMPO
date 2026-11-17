package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateActivationCode membuat kode unik sekali pakai format ACT-XXXX-XXXX (alfanumerik kapital)
func GenerateActivationCode() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Hindari karakter ambigu O/0, I/1
	part1 := make([]byte, 4)
	part2 := make([]byte, 4)

	for i := range part1 {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		part1[i] = charset[num.Int64()]
	}

	for i := range part2 {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		part2[i] = charset[num.Int64()]
	}

	return fmt.Sprintf("ACT-%s-%s", string(part1), string(part2)), nil
}
