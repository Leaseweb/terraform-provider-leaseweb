package shared

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/shared"
)

func TestNewSdkError(t *testing.T) {
	t.Run("expected Error is returned", func(t *testing.T) {
		err := errors.New("tralala")
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
			msg: "prefix: tralala",
			err: err,
			ErrorResponse: &shared.ErrorResponse{
				CorrelationId: "correlationId",
				ErrorCode:     "errorCode",
				ErrorMessage:  "errorMessage",
				ErrorDetails:  map[string][]string{"attribute": {"error1", "error2"}},
			},
		}

		assert.Equal(t, want, *got)
	})

	t.Run("invalid json does not return error", func(t *testing.T) {
		err := errors.New("tralala")
		response := http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
		}

		got := NewSdkError("prefix", err, &response)
		assert.Nil(t, got.ErrorResponse)
	})

	t.Run("nothing breaks when httpResponse is nil", func(t *testing.T) {
		err := errors.New("tralala")

		got := NewSdkError("prefix", err, nil)
		assert.Nil(t, got.ErrorResponse)
	})

}

func TestRepositoryError_Error(t *testing.T) {
	err := RepositoryError{msg: "tralala"}
	want := "tralala"
	got := err.Error()

	assert.Equal(t, want, got)
}

func TestNewGeneralError(t *testing.T) {
	err := errors.New("tralala")

	got := NewGeneralError("prefix", err)
	want := RepositoryError{
		msg: "prefix: tralala",
		err: err,
	}

	assert.Equal(t, want, *got)
}
