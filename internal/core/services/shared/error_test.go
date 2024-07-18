package shared

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	sharedRepository "terraform-provider-leaseweb/internal/repositories/shared"
)

func TestNewRepositoryError(t *testing.T) {
	err := errors.New("tralala")
	response := http.Response{StatusCode: 500, Body: nil}

	repositoryError := sharedRepository.NewSdkError(
		"repositoryErrorPrefix",
		err,
		&response,
	)
	got := NewRepositoryError("prefix", repositoryError)

	want := ServiceError{
		msg:             "prefix: repositoryErrorPrefix: tralala",
		RepositoryError: repositoryError,
	}

	assert.Equal(t, want, *got)
}

func TestServiceError_Error(t *testing.T) {
	err := ServiceError{msg: "tralala"}
	want := "tralala"
	got := err.Error()

	assert.Equal(t, want, got)
}

func TestNewGeneralError(t *testing.T) {
	err := errors.New("tralala")

	got := NewGeneralError("prefix", err)

	want := ServiceError{
		msg:          "prefix: tralala",
		GeneralError: err,
	}

	assert.Equal(t, want, *got)
}
