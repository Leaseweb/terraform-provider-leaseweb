package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstances_OrderById(t *testing.T) {
	instances := Instances{Instance{Id: "d"}, Instance{Id: "2"}, Instance{Id: "c"}}

	got := instances.OrderById()
	want := Instances{Instance{Id: "2"}, Instance{Id: "c"}, Instance{Id: "d"}}

	assert.Equal(t, want, got)
}
