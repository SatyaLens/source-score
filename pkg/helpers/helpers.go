package helpers

import (
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
