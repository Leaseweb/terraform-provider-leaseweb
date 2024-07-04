package modify_plan

import (
	"slices"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type Regions []string

func (r Regions) Contains(region string) bool {
	return slices.Contains(r, region)
}

func NewRegions(sdkRegions []publicCloud.Region) Regions {
	var regions Regions
	for _, region := range sdkRegions {
		regions = append(regions, region.GetName())
	}

	return regions
}
