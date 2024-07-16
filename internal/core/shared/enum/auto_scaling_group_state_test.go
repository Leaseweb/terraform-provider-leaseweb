package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoScalingGroupState_String(t *testing.T) {
	got := AutoScalingGroupStateScaling.String()
	want := "SCALING"

	assert.Equal(t, want, got)
}

func TestNewAutoScalingGroupState(t *testing.T) {
	want := AutoScalingGroupStateUpdating
	got, err := NewAutoScalingGroupState("UPDATING")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
