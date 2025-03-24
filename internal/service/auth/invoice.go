package auth

import (
	"math/rand"
)

const IDChars = "123456789"

func generateInvoiceID(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	invoiceIDCharsLength := len(IDChars)
	for i := 0; i < length; i++ {
		buffer[i] = IDChars[int(buffer[i])%invoiceIDCharsLength]
	}

	return string(buffer), nil
}
