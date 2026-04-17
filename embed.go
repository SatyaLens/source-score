package embed

import _ "embed"

var (
	//go:embed api/source-score.yaml
	OpenAPI []byte
)
