package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorResponse(t *testing.T) {
	t.Run("json is processed correctly", func(t *testing.T) {
		jsonStr := `
{
  "correlationId": "correlationId",
  "errorCode": "errorCode",
  "errorMessage": "errorMessage",
  "errorDetails":  {
    "attribute": ["error1", "error2"]
  }
}
`
		want := ErrorResponse{
			CorrelationId: "correlationId",
			ErrorCode:     "errorCode",
			ErrorMessage:  "errorMessage",
			ErrorDetails:  map[string][]string{"attribute": {"error1", "error2"}},
		}
		got, err := NewErrorResponse(jsonStr)

		assert.NoError(t, err)
		assert.Equal(t, want, *got)
	})

	t.Run("Invalid json returns error", func(t *testing.T) {
		_, err := NewErrorResponse("")
		assert.Error(t, err)
	})
}

func ExampleNewErrorResponse() {
	errorResponse, _ := NewErrorResponse(`
  {
    "correlationId": "correlationId",
    "errorCode": "errorCode",
    "errorMessage": "errorMessage",
    "errorDetails":  {
      "attribute": ["error1", "error2"]
    }
  }
`)

	fmt.Println(errorResponse)
	// Output: &{correlationId errorCode errorMessage map[attribute:[error1 error2]]}
}
