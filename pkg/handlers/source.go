package handlers

type SourceHandler struct {
}

func NewSourceHandler() *SourceHandler {
	return &SourceHandler{}
}

func (sh *SourceHandler) CreateSource() string {
	return "unimplemented"
}

func (sh *SourceHandler) DeleteSource() string {
	return "unimplemented"
}

func (sh *SourceHandler) GetSource() string {
	return "unimplemented"
}

func (sh *SourceHandler) UpdateSource() string {
	return "unimplemented"
}