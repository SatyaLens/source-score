package acceptance_test

import (
	"log"
	"net"
	"os"
	"source-score/pkg/helpers"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type responseBody struct {
	Data string `json:"data"`
}

var (
	baseUrl string

	serverPort = os.Getenv("PORT")
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
