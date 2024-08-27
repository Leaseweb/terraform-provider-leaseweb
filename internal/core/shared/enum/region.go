package enum

import (
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/shared/enum_utils"
)

type Region string

func (r Region) String() string {
	return string(r)
}

const (
	RegionEuWest3      = "eu-west-3"
	RegionUsEast1      = "us-east-1"
	RegionEuCentral1   = "eu-central-1"
	RegionApSoutheast1 = "api-southeast-1"
	RegionUsWest1      = "us-west-1"
	RegionEuWest2      = "eu-west-2"
	RegionCaCentral1   = "ca-central-1"
)

var regions = []Region{
	RegionEuWest3,
	RegionUsEast1,
	RegionEuCentral1,
	RegionApSoutheast1,
	RegionUsWest1,
	RegionEuWest2,
	RegionCaCentral1,
}

func NewRegion(r string) (Region, error) {
	return enum_utils.FindEnumForString(r, regions, RegionEuWest3)
}
