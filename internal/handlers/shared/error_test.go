package shared

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	sharedService "terraform-provider-leaseweb/internal/core/services/shared"
	sharedRepository "terraform-provider-leaseweb/internal/repositories/shared"
)

func TestNewServiceError(t *testing.T) {
	err := errors.New("tralala")
	repositoryError := sharedRepository.NewGeneralError(
		"repositoryErrorPrefix",
		err,
	)

	serviceError := sharedService.NewRepositoryError(
		"serviceErrorPrefix",
		repositoryError,
	)
	got := NewServiceError("prefix", serviceError)

	want := HandlerError{
		msg:          "prefix: serviceErrorPrefix: repositoryErrorPrefix: tralala",
		ServiceError: serviceError,
	}

	assert.Equal(t, want, *got)
}

func TestHandlerError_Error(t *testing.T) {
	err := HandlerError{msg: "tralala"}
	want := "tralala"
	got := err.Error()

	assert.Equal(t, want, got)
}

func TestNewGeneralError(t *testing.T) {
	err := errors.New("tralala")

	got := NewGeneralError("prefix", err)

	want := HandlerError{
		msg:          "prefix: tralala",
		GeneralError: err,
	}

	assert.Equal(t, want, *got)
}
