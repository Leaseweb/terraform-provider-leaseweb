package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoScalingGroupType_String(t *testing.T) {
	got := AutoScalingCpuTypeManual.String()

	assert.Equal(t, "MANUAL", got)
}
