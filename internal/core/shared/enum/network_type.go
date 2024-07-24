package enum

import (
	"terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type NetworkType string

func (n NetworkType) String() string {
	return string(n)
}

const (
	NetworkTypeInternal NetworkType = "INTERNAL"
	NetworkTypePublic   NetworkType = "PUBLIC"
)

var networkTypes = []NetworkType{NetworkTypeInternal, NetworkTypePublic}

func NewNetworkType(s string) (NetworkType, error) {
	return enum_utils.FindEnumForString(s, networkTypes, NetworkTypePublic)
}
