package claim_test

import (
	"source-score/pkg/api"
	"source-score/pkg/domain/claim"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

const (
	testDBFile = "test_claims_database.db"
)

var (
	err          error
	sampleClaim1 api.Claim
	sampleClaim2 api.Claim
	claimRepo    claim.ClaimRepository
	testDB       *gorm.DB
)

func TestClaim(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Claim Unit Test Suite")
}
