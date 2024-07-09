package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethod_String(t *testing.T) {
	got := MethodGet.String()

	assert.Equal(t, "GET", got)

}
