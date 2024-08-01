package shared

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSdkError(t *testing.T) {
	t.Run("expected Error is returned", func(t *testing.T) {
		err := errors.New("result")
		response := http.Response{
			StatusCode: 500,
			Body: io.NopCloser(
				bytes.NewReader(
					[]byte(`
{
  "correlationId": "correlationId",
  "errorCode": "errorCode",
  "errorMessage": "errorMessage",
  "errorDetails":  {
    "attribute": ["error1", "error2"]
  }
}
          `),
				),
			),
		}

		got := NewSdkError("prefix", err, &response)
		want := RepositoryError{
			msg: "prefix: result",
			err: err,
			ErrorResponse: &ErrorResponse{
				CorrelationId: "correlationId",
				ErrorCode:     "errorCode",
				ErrorMessage:  "errorMessage",
				ErrorDetails:  map[string][]string{"attribute": {"error1", "error2"}},
			},
		}

		assert.Equal(t, want, *got)
	})

	t.Run("invalid json does not return error", func(t *testing.T) {
		err := errors.New("Result")
		response := http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
		}

		got := NewSdkError("prefix", err, &response)
		assert.Nil(t, got.ErrorResponse)
	})

	t.Run("nothing breaks when httpResponse is nil", func(t *testing.T) {
		err := errors.New("Result")

		got := NewSdkError("prefix", err, nil)
		assert.Nil(t, got.ErrorResponse)
	})

}

func TestRepositoryError_Error(t *testing.T) {
	err := RepositoryError{msg: "Result"}
	want := "Result"
	got := err.Error()

	assert.Equal(t, want, got)
}

func TestNewGeneralError(t *testing.T) {
	err := errors.New("Result")

	got := NewGeneralError("prefix", err)
	want := RepositoryError{
		msg: "prefix: Result",
		err: err,
	}

	assert.Equal(t, want, *got)
}

func ExampleNewSdkError() {
	httpResponse := http.Response{
		StatusCode: 500,
		Body: io.NopCloser(
			bytes.NewReader(
				[]byte(`
{
  "correlationId": "correlationId",
  "errorCode": "errorCode",
  "errorMessage": "errorMessage",
  "errorDetails":  {
    "attribute": ["error1", "error2"]
  }
}
          `),
			),
		),
	}

	repositoryError := NewSdkError(
		"prefix",
		errors.New("some error"),
		&httpResponse,
	)

	fmt.Println(repositoryError)
	fmt.Println(repositoryError.ErrorResponse)
	// Output:
	// prefix: some error
	// &{correlationId errorCode errorMessage map[attribute:[error1 error2]]}
}

func ExampleNewGeneralError() {
	repositoryError := NewGeneralError(
		"prefix",
		errors.New("some error"),
	)

	fmt.Println(repositoryError)
	// Output: prefix: some error
}
