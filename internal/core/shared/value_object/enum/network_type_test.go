package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkType_String(t *testing.T) {
	got := NetworkTypeInternal.String()

	assert.Equal(t, "INTERNAL", got)

}
