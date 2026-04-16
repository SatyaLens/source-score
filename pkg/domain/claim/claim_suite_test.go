package claim_test

import (
	"context"
	"source-score/pkg/api"
	"source-score/pkg/db/pgsql"
	"source-score/pkg/domain/claim"
	"source-score/pkg/helpers"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	claim1Digest = "369d9f3047c66c2e9b5e39693d9de3664b61a36a2d77cd0484fade042350d4a1"
	claim2Digest = "a96fe15d3040685b06d0c195d54a13692a9002db148498f185babfb6a083f801"
	srcDigest    = "8649a4126fb4fc9a750f432b729c8477398cf28ca241403b2cd36a6dc841f441"
	testDBFile   = "test_claims_database.db"
)

var (
	err       error
	claimRepo claim.ClaimRepository
	testDB    *gorm.DB

	// create a source to reference from claims
	sampleSource = api.Source{
		Name:      "Sample Source 1",
		Score:     0,
		Summary:   "Sample source summary",
		Tags:      "tag1",
		Uri:       "https://sample-uri-1",
		UriDigest: srcDigest,
	}

	sampleClaim1 = api.Claim{
		SourceUriDigest: sampleSource.UriDigest,
		Summary:         "Sample claim summary 1",
		Title:           "Sample Claim 1",
		Uri:             "https://sample-claim-1",
		UriDigest:       claim1Digest,
		Checked:         false,
		Validity:        false,
	}

	sampleClaim2 = api.Claim{
		SourceUriDigest: sampleSource.UriDigest,
		Summary:         "Sample claim summary 2",
		Title:           "Sample Claim 2",
		Uri:             "https://sample-claim-2",
		UriDigest:       claim2Digest,
		Checked:         false,
		Validity:        false,
	}
)

func TestClaim(t *testing.T) {
	var _ = BeforeSuite(func() {
		testDB, err = gorm.Open(sqlite.Open(testDBFile))
		Expect(err).ToNot(HaveOccurred())

		err = testDB.AutoMigrate(&api.Source{}, &api.Claim{}, &api.Proof{})
		Expect(err).ToNot(HaveOccurred())

		result := testDB.Create(&sampleSource)
		Expect(result.Error).ToNot(HaveOccurred())

		// configure repository
		claimRepo = claim.NewClaimRepository(context.TODO(), &pgsql.Client{DB: testDB})
	})

	var _ = AfterSuite(func() {
		db, err := testDB.DB()
		Expect(err).ToNot(HaveOccurred())
		db.Close()
		err = helpers.DeleteFileIfExists(testDBFile)
		Expect(err).ToNot(HaveOccurred())
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Claim Unit Test Suite")
}
