package public_cloud_repository

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPublicCloudRepositoryError(t *testing.T) {
	err := errors.New("tralala")
	response := http.Response{StatusCode: 500, Body: nil}

	got := newPublicCloudRepositoryError("prefix", err, &response)
	want := PublicCloudRepositoryError{
		msg:             "prefix: tralala",
		err:             err,
		SdkHttpResponse: &response,
	}

	assert.Equal(t, want, got)
}

func TestPublicCloudRepositoryError_Error(t *testing.T) {
	err := PublicCloudRepositoryError{msg: "tralala"}
	want := "tralala"
	got := err.Error()

	assert.Equal(t, want, got)
}
