package source

type SourceService struct {
	sourceRepo *sourceRepository
}

func NewSourceService(sourceRepo *sourceRepository) *SourceService {
	return &SourceService{
		sourceRepo: sourceRepo,
	}
}

