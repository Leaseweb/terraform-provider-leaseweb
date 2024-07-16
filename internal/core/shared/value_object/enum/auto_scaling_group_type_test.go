package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoScalingGroupType_String(t *testing.T) {
	got := AutoScalingCpuTypeManual.String()

	assert.Equal(t, "MANUAL", got)
}

func TestNewAutoScalingGroupType(t *testing.T) {
	want := AutoScalingGroupTypeCpuBased
	got, err := NewAutoScalingGroupType("CPU_BASED")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
