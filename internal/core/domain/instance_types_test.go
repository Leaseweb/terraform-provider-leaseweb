package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceTypes_ToArray(t *testing.T) {
	instanceTypes := InstanceTypes{InstanceType{Name: "tralala"}}
	got := instanceTypes.ToArray()
	want := []string{"tralala"}

	assert.Equal(t, want, got)
}
