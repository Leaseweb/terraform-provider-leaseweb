package shared

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	sharedService "terraform-provider-leaseweb/internal/core/services/shared"
	sharedRepository "terraform-provider-leaseweb/internal/repositories/shared"
)

func TestNewFromServiceError(t *testing.T) {
	err := errors.New("tralala")
	errorResponse := sharedRepository.ErrorResponse{ErrorCode: "123"}

	repositoryError := sharedRepository.NewGeneralError(
		"repositoryErrorPrefix",
		err,
	)
	repositoryError.ErrorResponse = &errorResponse

	serviceError := sharedService.NewFromRepositoryError(
		"serviceErrorPrefix",
		*repositoryError,
	)

	got := NewFromServicesError("prefix", serviceError)

	want := HandlerError{
		msg:           "prefix: serviceErrorPrefix: repositoryErrorPrefix: tralala",
		ErrorResponse: &errorResponse,
	}

	assert.Equal(t, want, *got)
}

func TestHandlerError_Error(t *testing.T) {
	err := HandlerError{msg: "tralala"}
	want := "tralala"
	got := err.Error()

	assert.Equal(t, want, got)
}

func TestNewError(t *testing.T) {
	err := errors.New("tralala")

	got := NewError("prefix", err)

	want := HandlerError{
		msg:          "prefix: tralala",
		GeneralError: err,
	}

	assert.Equal(t, want, *got)
}

func ExampleNewFromServicesError() {
	handlerError := NewFromServicesError(
		"handlerPrefix",
		sharedService.NewError("sharedPrefix", errors.New("tralala")),
	)

	fmt.Println(handlerError)
	// Output: handlerPrefix: sharedPrefix: tralala
}

func ExampleNewError() {
	handlerError := NewError("prefix", errors.New("tralala"))

	fmt.Println(handlerError)
	// Output: prefix: tralala
}
