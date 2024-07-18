package shared

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepositoryError(t *testing.T) {
	err := errors.New("tralala")
	response := http.Response{StatusCode: 500, Body: nil}

	got := NewRepositoryError("prefix", err, &response)
	want := RepositoryError{
		msg:             "prefix: tralala",
		err:             err,
		SdkHttpResponse: &response,
	}

	assert.Equal(t, want, got)
}

func TestRepositoryError_Error(t *testing.T) {
	err := RepositoryError{msg: "tralala"}
	want := "tralala"
	got := err.Error()

	assert.Equal(t, want, got)
}
