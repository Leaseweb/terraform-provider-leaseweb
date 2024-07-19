package shared

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	sharedRepository "terraform-provider-leaseweb/internal/repositories/shared"
	"terraform-provider-leaseweb/internal/shared"
)

func TestNewFromRepositoryError(t *testing.T) {
	err := errors.New("tralala")
	response := http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewReader([]byte(""))),
	}

	errorResponse := shared.ErrorResponse{ErrorCode: "54"}

	repositoryError := sharedRepository.NewSdkError(
		"repositoryErrorPrefix",
		err,
		&response,
	)
	repositoryError.ErrorResponse = &errorResponse

	got := NewFromRepositoryError("prefix", repositoryError)

	want := ServiceError{
		msg:           "prefix: repositoryErrorPrefix: tralala",
		ErrorResponse: &errorResponse,
	}

	assert.Equal(t, want, *got)
}

func TestServiceError_Error(t *testing.T) {
	err := ServiceError{msg: "tralala"}
	want := "tralala"
	got := err.Error()

	assert.Equal(t, want, got)
}

func TestNewError(t *testing.T) {
	err := errors.New("tralala")

	got := NewError("prefix", err)

	want := ServiceError{
		msg:          "prefix: tralala",
		GeneralError: err,
	}

	assert.Equal(t, want, *got)
}
