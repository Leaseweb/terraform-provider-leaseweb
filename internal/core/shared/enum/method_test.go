package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethod_String(t *testing.T) {
	got := MethodGet.String()

	assert.Equal(t, "GET", got)

}

func TestNewMethod(t *testing.T) {
	want := MethodHead
	got, err := NewMethod("HEAD")

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
