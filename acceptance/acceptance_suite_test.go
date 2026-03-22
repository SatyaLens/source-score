package acceptance_test

import (
	"log"
	"net"
	"net/http"
	"os"
	"source-score/pkg/api"
	"source-score/pkg/helpers"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type responseBody struct {
	Data string `json:"data"`
}

const (
	uriDigest1 = "8649a4126fb4fc9a750f432b729c8477398cf28ca241403b2cd36a6dc841f441"
)

var (
	baseUrl string

	client     = &http.Client{Timeout: 10 * time.Second}
	serverPort = os.Getenv("PORT")
	sourceInput1 = api.SourceInput{
		Name:    "Sample Source 1",
		Summary: "Sample summary",
		Tags:    "tag1",
		Uri:     "https://sample-uri-1",
	}
)

func TestSourceScore(t *testing.T) {
	if serverPort == "" {
		serverPort = "8080"
	}

	if !isLocalPortOpen(serverPort) {
		log.Fatalf("application not running on port: %s", serverPort)
	}

	baseUrl = "http://" + helpers.Localhost + ":" + serverPort

	RegisterFailHandler(Fail)
	RunSpecs(t, "SourceScore Acceptance Test Suite")
}

func isLocalPortOpen(port string) bool {
	address := net.JoinHostPort(helpers.Localhost, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
