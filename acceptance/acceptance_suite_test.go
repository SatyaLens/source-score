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
	uriDigest1   = "8649a4126fb4fc9a750f432b729c8477398cf28ca241403b2cd36a6dc841f441"
	uriDigest2   = "978d81ca657062910f60263c26ce7261e7530e53bfd28aa48748155eb5621868"
	uriDigest3   = "f70fe06de54dcaa05e3fcda03ae724ad9d8603c04f6cdbd838c0ad4f2e789ba1"
	claim1Digest = "369d9f3047c66c2e9b5e39693d9de3664b61a36a2d77cd0484fade042350d4a1"
	claim2Digest = "a96fe15d3040685b06d0c195d54a13692a9002db148498f185babfb6a083f801"
)

var (
	baseUrl string

	client       = &http.Client{Timeout: 10 * time.Second}
	serverPort   = os.Getenv("PORT")
	sourceInput1 = api.SourceInput{
		Name:    "Sample Source 1",
		Summary: "Sample summary",
		Tags:    "tag1",
		Uri:     "https://sample-uri-1",
	}
	sourceInput2 = api.SourceInput{
		Name:    "Sample Source 2",
		Summary: "Sample summary 2",
		Tags:    "tag2",
		Uri:     "https://sample-uri-2",
	}
	sourceInput3 = api.SourceInput{
		Name:    "Sample Source 3",
		Summary: "Sample summary 3",
		Tags:    "tag2",
		Uri:     "https://sample-uri-3",
	}
	sampleSource1 = api.Source{
		Name:      "Sample Source 1",
		Score:     0,
		Summary:   "Sample summary",
		Tags:      "tag1",
		Uri:       "https://sample-uri-1",
		UriDigest: uriDigest1,
	}
	sampleSource2 = api.Source{
		Name:      "Sample Source 2",
		Score:     0,
		Summary:   "Sample summary 2",
		Tags:      "tag2",
		Uri:       "https://sample-uri-2",
		UriDigest: uriDigest2,
	}
	sampleClaim1 = api.Claim{
		SourceUriDigest: uriDigest3,
		Summary:         "Sample claim summary 1",
		Title:           "Sample Claim 1",
		Uri:             "https://sample-claim-1",
		UriDigest:       claim1Digest,
		Checked:         false,
		Validity:        false,
	}
	sampleClaim2 = api.Claim{
		SourceUriDigest: uriDigest3,
		Summary:         "Sample claim summary 2",
		Title:           "Sample Claim 2",
		Uri:             "https://sample-claim-2",
		UriDigest:       claim2Digest,
		Checked:         false,
		Validity:        false,
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
