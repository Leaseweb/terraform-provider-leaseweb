package enum

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

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
	return enum_utils.FindEnumForString(value, methods, MethodGet)
}
