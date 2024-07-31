package shared

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	sharedService "terraform-provider-leaseweb/internal/core/services/errors"
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

	want := FacadeError{
		msg:           "prefix: serviceErrorPrefix: repositoryErrorPrefix: tralala",
		ErrorResponse: &errorResponse,
	}

	assert.Equal(t, want, *got)
}

func TestFacadeError_Error(t *testing.T) {
	err := FacadeError{msg: "tralala"}
	want := "tralala"
	got := err.Error()

	assert.Equal(t, want, got)
}

func TestNewError(t *testing.T) {
	err := errors.New("tralala")

	got := NewError("prefix", err)

	want := FacadeError{
		msg:          "prefix: tralala",
		GeneralError: err,
	}

	assert.Equal(t, want, *got)
}

func ExampleNewFromServicesError() {
	facadeError := NewFromServicesError(
		"facadePrefix",
		sharedService.NewError("sharedPrefix", errors.New("tralala")),
	)

	fmt.Println(facadeError)
	// Output: facadePrefix: sharedPrefix: tralala
}

func ExampleNewError() {
	facadeError := NewError("prefix", errors.New("tralala"))

	fmt.Println(facadeError)
	// Output: prefix: tralala
}
