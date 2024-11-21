package acceptance

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	baseUrl string

	port = os.Getenv("PORT")
)

func TestSourceScore(t *testing.T) {
	if port == "" {
		port = "8080"
	}
	baseUrl = "localhost:" + port

	RegisterFailHandler(Fail)
	RunSpecs(t, "SourceScore Acceptance Test Suite")
}
