package utils

import (
	"crypto/rand"
	"io"
)

// GenerateCryptoOTP menghasilkan string digit angka acak kriptografis dengan panjang tertentu.
// Fungsi ini memanfaatkan CSPRNG (Cryptographically Secure Pseudo-Random Number Generator)
// dari kernel sistem operasi (/dev/urandom) untuk menjamin bahwa kode tidak dapat diprediksi,
// serta mempertahankan zero-padding (misal "012345" tetap 6 digit utuh).
func GenerateCryptoOTP(length int) (string, error) {
	const digits = "0123456789"
	bytes := make([]byte, length)

	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = digits[bytes[i]%byte(len(digits))]
	}

	return string(bytes), nil
}
