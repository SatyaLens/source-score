package helpers

const (
	Localhost       = "127.0.0.1"
	RequestIdHeader = "X-Request-ID"
)

var (
	ApiReqLogFields = []string {
		RequestIdHeader,
	}
)