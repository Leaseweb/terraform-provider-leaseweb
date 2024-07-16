package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkType_String(t *testing.T) {
	got := NetworkTypeInternal.String()

	assert.Equal(t, "INTERNAL", got)

}

func TestNewNetworkType(t *testing.T) {
	want := NetworkTypeInternal
	got, err := NewNetworkType("INTERNAL")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
