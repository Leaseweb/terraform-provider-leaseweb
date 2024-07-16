package enum

type Method string

func (m Method) String() string {
	return string(m)
}

const (
	MethodGet     Method = "GET"
	MethodHead    Method = "HEAD"
	MethodPost    Method = "POST"
	MethodOptions Method = "OPTIONS"
)

var methods = []Method{MethodGet, MethodHead, MethodPost, MethodOptions}

func NewMethod(value string) (Method, error) {
	return findEnumForString(value, methods, MethodGet)
}
