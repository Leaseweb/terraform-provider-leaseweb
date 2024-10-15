package errors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	repository2 "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewFromRepositoryError(t *testing.T) {
	err := errors.New("tralala")
	response := http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewReader([]byte(""))),
	}

	errorResponse := repository2.ErrorResponse{ErrorCode: "54"}

	repositoryError := repository2.NewSdkError(
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
	repositoryError := repository2.NewSdkError(
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
