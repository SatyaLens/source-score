package claim_test

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "log"

    "source-score/pkg/api"
    "source-score/pkg/db/pgsql"
    "source-score/pkg/domain/claim"
    "source-score/pkg/helpers"

    "testing"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestClaimRepository(t *testing.T) {
    var _ = BeforeSuite(func() {
        testDB, err = gorm.Open(sqlite.Open(testDBFile))
        Expect(err).ToNot(HaveOccurred())

        err = testDB.AutoMigrate(&api.Source{}, &api.Claim{})
        Expect(err).ToNot(HaveOccurred())

        // create a source to reference from claims
        sampleSource := api.Source{
            Name:      "Sample Source 1",
            Score:     0,
            Summary:   "Sample source summary",
            Tags:      "tag1",
            Uri:       "https://sample-uri-1",
            UriDigest: "source-digest-1",
        }

        result := testDB.Create(&sampleSource)
        Expect(result.Error).ToNot(HaveOccurred())

        // create two sample claims
        claimUri1 := "https://sample-claim-1"
        h1 := sha256.New()
        _, err = h1.Write([]byte(claimUri1))
        Expect(err).ToNot(HaveOccurred())
        digest1 := hex.EncodeToString(h1.Sum(nil))

        claimUri2 := "https://sample-claim-2"
        h2 := sha256.New()
        _, err = h2.Write([]byte(claimUri2))
        Expect(err).ToNot(HaveOccurred())
        digest2 := hex.EncodeToString(h2.Sum(nil))

        sampleClaim1 = api.Claim{
            SourceUriDigest: sampleSource.UriDigest,
            Summary:         "Sample claim summary 1",
            Title:           "Sample Claim 1",
            Uri:             claimUri1,
            UriDigest:       digest1,
            Checked:         false,
            Validity:        false,
        }

        sampleClaim2 = api.Claim{
            SourceUriDigest: sampleSource.UriDigest,
            Summary:         "Sample claim summary 2",
            Title:           "Sample Claim 2",
            Uri:             claimUri2,
            UriDigest:       digest2,
            Checked:         true,
            Validity:        true,
        }

        result = testDB.Create(&sampleClaim1)
        Expect(result.Error).ToNot(HaveOccurred())

        result = testDB.Create(&sampleClaim2)
        Expect(result.Error).ToNot(HaveOccurred())

        // configure repository
        claimRepo = claim.NewClaimRepository(context.TODO(), &pgsql.Client{DB: testDB})
    })

    var _ = AfterSuite(func() {
        log.Println("deleting test files")
        err = helpers.DeleteFileIfExists(testDBFile)
        Expect(err).ToNot(HaveOccurred())
    })

    RegisterFailHandler(Fail)
    RunSpecs(t, "Claim Repository Unit Test Suite")
}

var _ = Describe("Claim repository layer unit tests", Ordered, func() {
    Context("Happy path", Ordered, func() {
        When("Retrieving all claims from the DB", func() {
            It("Should return all claim records from the DB", func() {
                claims, err := claimRepo.GetClaims(context.TODO())
                Expect(err).ToNot(HaveOccurred())
                Expect(len(claims)).To(Equal(2))

                Expect(claims).To(ContainElements(
                    sampleClaim1,
                    sampleClaim2,
                ))
            })
        })
    })
})
package claim_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)
