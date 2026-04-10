package proof_test

import (
	"context"
	"source-score/pkg/api"
	"source-score/pkg/db/pgsql"
	"source-score/pkg/domain/proof"
	"source-score/pkg/helpers"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	proof1Digest = "6f2479b5249b1c27c4935da5594bc72bb0b9e59e704aea9af50780bc6178c357"
	proof2Digest = "8df5229f310ae8322062834f3ba45a38ecef8ded549665d1170e15c8249b7cd0"
	claimDigest  = "369d9f3047c66c2e9b5e39693d9de3664b61a36a2d77cd0484fade042350d4a1"
	testDBFile   = "test_proofs_database.db"
)

var (
	err       error
	proofRepo proof.ProofRepository
	testDB    *gorm.DB

	sampleClaim = api.Claim{
		SourceUriDigest: "8649a4126fb4fc9a750f432b729c8477398cf28ca241403b2cd36a6dc841f441",
		Summary:         "Sample claim for proof tests",
		Title:           "Sample Claim",
		Uri:             "https://sample-claim",
		UriDigest:       claimDigest,
		Checked:         false,
		Validity:        false,
	}

	sampleProof1 = api.Proof{
		ClaimUriDigest: sampleClaim.UriDigest,
		ReviewedBy:     "ReviewerA",
		SupportsClaim:  true,
		Uri:            "https://sample-proof-1",
		UriDigest:      proof1Digest,
	}

	sampleProof2 = api.Proof{
		ClaimUriDigest: sampleClaim.UriDigest,
		ReviewedBy:     "ReviewerB",
		SupportsClaim:  false,
		Uri:            "https://sample-proof-2",
		UriDigest:      proof2Digest,
	}
)

func TestProof(t *testing.T) {
	var _ = BeforeSuite(func() {
		testDB, err = gorm.Open(sqlite.Open(testDBFile))
		Expect(err).ToNot(HaveOccurred())

		err = testDB.AutoMigrate(&api.Source{}, &api.Claim{}, &api.Proof{})
		Expect(err).ToNot(HaveOccurred())

		// insert claim to reference from proofs
		result := testDB.Create(&sampleClaim)
		Expect(result.Error).ToNot(HaveOccurred())

		// configure repository
		proofRepo = proof.NewProofRepository(context.TODO(), &pgsql.Client{DB: testDB})
	})

	var _ = AfterSuite(func() {
		db, err := testDB.DB()
		Expect(err).ToNot(HaveOccurred())
		db.Close()
		err = helpers.DeleteFileIfExists(testDBFile)
		Expect(err).ToNot(HaveOccurred())
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Proof Unit Test Suite")
}
