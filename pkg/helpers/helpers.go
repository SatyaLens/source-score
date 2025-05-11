package helpers

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
)

func DeleteFileIfExists(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Printf("file %s does not exist\n", filePath)
		return nil
	}

	return os.Remove(filePath)
}

func GetSHA256Hash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}
