package errors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/repository/shared"
	"github.com/stretchr/testify/assert"
)

func TestNewFromRepositoryError(t *testing.T) {
	err := errors.New("tralala")
	response := http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewReader([]byte(""))),
	}

	errorResponse := shared.ErrorResponse{ErrorCode: "54"}

	repositoryError := shared.NewSdkError(
		"repositoryErrorPrefix",
		err,
		&response,
	)
	repositoryError.ErrorResponse = &errorResponse

	got := NewFromRepositoryError("prefix", *repositoryError)

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

func ExampleNewFromRepositoryError() {
	repositoryError := shared.NewSdkError(
		"repositoryErrorPrefix",
		errors.New("sdk error"),
		nil,
	)

	fromRepositoryError := NewFromRepositoryError(
		"prefix",
		*repositoryError,
	)

	fmt.Println(fromRepositoryError)
	// Output: prefix: repositoryErrorPrefix: sdk error
}

func ExampleNewError() {
	newError := NewError("prefix", errors.New("tralala"))

	fmt.Println(newError)
	// Output: prefix: tralala
}
