package helpers

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

func DeleteFileIfExists(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Printf("file %s does not exist\n", filePath)
		return nil
	}

	return os.Remove(filePath)
}

func ValidateNonEmpty(fl validator.FieldLevel) bool {
	return fl.Field().String() != ""
}

func ValidateNoSpace(fl validator.FieldLevel) bool {
	return !strings.ContainsAny(fl.Field().String(), " \t\n\r")
}

func ValidateHttpsURL(fl validator.FieldLevel) bool {
	raw := fl.Field().String()
	u, err := url.ParseRequestURI(raw)
	return err == nil && u.Scheme == "https" && u.Host != ""
}

func SameHost(url1, url2 string) (bool, error) {
	a, err := url.Parse(url1)
	if err != nil {
		return false, fmt.Errorf("invalid URL %q: %w", url1, err)
	}
	b, err := url.Parse(url2)
	if err != nil {
		return false, fmt.Errorf("invalid URL %q: %w", url2, err)
	}
	return a.Hostname() == b.Hostname(), nil
}
