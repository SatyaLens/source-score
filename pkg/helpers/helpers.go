package helpers

import (
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

// TODO: add uri digest validation here
func ValidateUriDigest(uriDigest string) error {
	if uriDigest == "" {
		return fmt.Errorf("invalid uri digest")
	}
	return nil
}
