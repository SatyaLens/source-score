package source

type sourceRepository struct {
	client interface{}
}

func NewSourceRepository(client interface{}) *sourceRepository {
	return &sourceRepository{
		client: client,
	}
}

func (sr *sourceRepository) GetSourceByDigest(uriDigest string) {
	
}