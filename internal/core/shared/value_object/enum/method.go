package enum

type Method string

func (m Method) String() string {
	return string(m)
}

type Methods []Method

const (
	MethodGet     Method = "GET"
	MethodHead    Method = "HEAD"
	MethodPost    Method = "POST"
	MethodOptions Method = "OPTIONS"
)

var MethodValues = Methods{MethodGet, MethodHead, MethodPost, MethodOptions}
