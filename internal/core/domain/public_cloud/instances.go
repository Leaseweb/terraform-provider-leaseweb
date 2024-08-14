package public_cloud

import (
	"sort"
)

type Instances []Instance

func (inst Instances) OrderById() Instances {
	sort.Slice(inst, func(i, j int) bool {
		return inst[i].Id < inst[j].Id
	})

	return inst
}
