package telegramauth

import "crypto/sha256"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func getUserToken(input string, size int) string {
	hash := sha256.Sum256([]byte(input))

	result := make([]byte, size)

	for i := 0; i < size; i++ {
		byteValue := hash[i%len(hash)]
		result[i] = charset[byteValue%byte(len(charset))]
	}

	return string(result)
}
