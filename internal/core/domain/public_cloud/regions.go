package public_cloud

import (
	"fmt"
)

type ErrCannotFindRegion struct {
	msg string
}

func (e ErrCannotFindRegion) Error() string {
	return e.msg
}

type Regions []Region

func (r Regions) Contains(region string) bool {
	for _, r := range r {
		if r.Name == region {
			return true
		}
	}

	return false
}

func (r Regions) GetByName(region string) (*Region, error) {
	for _, r := range r {
		if r.Name == region {
			return &r, nil
		}
	}

	return nil, ErrCannotFindRegion{
		msg: fmt.Sprintf("region %q not found", region),
	}
}

func (r Regions) ToArray() []string {
	var values []string
	for _, region := range r {
		values = append(values, region.String())
	}

	return values
}
