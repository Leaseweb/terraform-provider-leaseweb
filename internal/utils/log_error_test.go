package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
)

func TestLogError(t *testing.T) {
	t.Run("response is set", func(t *testing.T) {

		diags := diag.Diagnostics{}

		LogError(
			context.TODO(),
			&ErrorResponse{},
			&diags,
			"summary",
			"detail",
		)

		assert.Equal(
			t,
			"summary",
			diags[0].Summary(),
			"error contains summary",
		)
		assert.Equal(
			t,
			"detail",
			diags[0].Detail(),
			"error contains detail",
		)
	})

	t.Run("response is not set", func(t *testing.T) {
		diags := diag.Diagnostics{}
		summary := "summary"
		detail := "detail"

		LogError(context.TODO(), nil, &diags, summary, detail)

		assert.Equal(
			t,
			"summary",
			diags[0].Summary(),
			"error contains summary",
		)
		assert.Equal(
			t,
			"detail",
			diags[0].Detail(),
			"error contains detail",
		)
	})
}

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
		want := SdkError{
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
		err := errors.New("result")
		response := http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
		}

		got := NewSdkError("prefix", err, &response)
		assert.Nil(t, got.ErrorResponse)
	})

	t.Run("nothing breaks when httpResponse is nil", func(t *testing.T) {
		err := errors.New("result")

		got := NewSdkError("prefix", err, nil)
		assert.Nil(t, got.ErrorResponse)
	})

}

func TestRepositoryError_Error(t *testing.T) {
	err := SdkError{msg: "Result"}
	want := "Result"
	got := err.Error()

	assert.Equal(t, want, got)
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

func Test_newErrorResponse(t *testing.T) {
	t.Run("response is processed correctly", func(t *testing.T) {
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

		want := ErrorResponse{
			CorrelationId: "correlationId",
			ErrorCode:     "errorCode",
			ErrorMessage:  "errorMessage",
			ErrorDetails:  map[string][]string{"attribute": {"error1", "error2"}},
		}
		got, err := newErrorResponse(httpResponse.Body)

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("Invalid json returns error", func(t *testing.T) {
		httpResponse := http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}

		_, err := newErrorResponse(httpResponse.Body)
		assert.Error(t, err)
	})
}
