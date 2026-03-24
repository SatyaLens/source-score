package source_test

import (
	"context"
	"log"
	"source-score/pkg/api"
	"source-score/pkg/db/pgsql"
	"source-score/pkg/domain/source"
	"source-score/pkg/domain/source/sourcefakes"
	"source-score/pkg/helpers"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	uriDigest1 = "8649a4126fb4fc9a750f432b729c8477398cf28ca241403b2cd36a6dc841f441"
	uriDigest2 = "978d81ca657062910f60263c26ce7261e7530e53bfd28aa48748155eb5621868"
	testDBFile = "test_database.db"
)

var (
	err                error
	sampleSource1      api.Source
	sampleSource2      api.Source
	sampleSourceInput1 api.SourceInput
	sampleSourceInput2 api.SourceInput
	sourceRepo         source.SourceRepository
	sourceSvc          source.SourceService
	testDB             *gorm.DB
	updatedSource      api.Source

	// fakes
	fakeSourceRepo = sourcefakes.FakeSourceRepository{}
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

		sampleSourceInput2 = api.SourceInput{
			Name:    "Sample Source 2",
			Summary: "Sample summary 2",
			Tags:    "tag2",
			Uri:     "https://sample-uri-2",
		}

		sampleSource1 = api.Source{
			Name:      "Sample Source 1",
			Score:     0.5,
			Summary:   "Sample summary",
			Tags:      "tag1",
			Uri:       "https://sample-uri-1",
			UriDigest: uriDigest1,
		}

		sampleSource2 = api.Source{
			Name:      "Sample Source 2",
			Score:     0.7,
			Summary:   "Sample summary 2",
			Tags:      "tag2",
			Uri:       "https://sample-uri-2",
			UriDigest: uriDigest2,
		}

		// configure layers
		sourceRepo = source.NewSourceRepository(context.TODO(), &pgsql.Client{
			DB: testDB,
		})
		sourceSvc = source.NewSourceService(context.TODO(), &fakeSourceRepo)
	})

	var _ = AfterSuite(func() {
		log.Println("deleting test files")
		err = helpers.DeleteFileIfExists(testDBFile)
		Expect(err).ToNot(HaveOccurred())
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Source Unit Test Suite")
}
