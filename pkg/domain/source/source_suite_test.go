package source_test

import (
	"context"
	"log"
	"source-score/pkg/api"
	"source-score/pkg/db/cnpg"
	"source-score/pkg/domain/source"
	"source-score/pkg/helpers"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	testDBFile = "test_database.db"
)

var (
	err                error
	sampleSourceInput1 api.SourceInput
	sourceRepo         *source.SourceRepository
	testDB             *gorm.DB
	uriDigest1         string
)

func TestSource(t *testing.T) {
	var _ = BeforeSuite(func() {
		testDB, err = gorm.Open(sqlite.Open(testDBFile))
		Expect(err).ToNot(HaveOccurred())

		err = testDB.AutoMigrate(&api.Source{}, &api.SourceInput{})
		Expect(err).ToNot(HaveOccurred())

		sampleSourceInput1 = api.SourceInput{
			Name:    "Sample Source 1",
			Summary: "Sample summary",
			Tags:    "tag1",
			Uri:     "https://sample-uri-1",
		}

		uriDigest1 = helpers.GetSHA256Hash(sampleSourceInput1.Uri)

		sourceRepo = source.NewSourceRepository(context.TODO(), &cnpg.Client{
			DB: testDB,
		})
	})

	var _ = AfterSuite(func() {
		log.Println("deleting test files")
		err = helpers.DeleteFileIfExists(testDBFile)
		Expect(err).ToNot(HaveOccurred())
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Source Unit Test Suite")
}
