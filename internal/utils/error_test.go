package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/stretchr/testify/assert"
)

func TestGetHttpErrorMessage(t *testing.T) {
	t.Run("error must return if response is nil", func(t *testing.T) {
		message := NewError(nil, errors.New("error content")).Error()
		assert.Equal(t, "error content", message)
	})

	t.Run("error must return if body of response is nil", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 500,
			Body:       nil,
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "error content", message)
	})

	t.Run("error must return if response is 2xx or 3xx", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("body message")),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "error content", message)
	})

	t.Run("error must return if response is not contain errorMessage", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString("body message")),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "error content", message)
	})

	t.Run("error must return if response is not contain errorMessage", func(t *testing.T) {
		errorContent, _ := json.Marshal(map[string]string{"balh": "balh"})
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(string(errorContent))),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "error content", message)
	})

	t.Run("errorMessage must return if response contains string errorMessage", func(t *testing.T) {
		errorContent, _ := json.Marshal(map[string]string{"errorMessage": "this is error message"})
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(string(errorContent))),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "this is error message", message)
	})

	t.Run("errorMessage must return if response contains object of errorMessage", func(t *testing.T) {
		errorContent, _ := json.Marshal(map[string]interface{}{"errorMessage": map[string]string{"a": "b", "c": "d"}})
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(string(errorContent))),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "map[a:b c:d]", message)
	})

	t.Run("errorMessage must return if response with error details if they exists", func(t *testing.T) {
		errorContent, _ := json.Marshal(map[string]interface{}{
			"errorMessage": map[string]string{"a": "b", "c": "d"},
			"errorDetails": map[string][]string{
				"password": {
					"this value should not be blank",
					"blah blah",
				},
				"email": {
					"this value should be valid",
					"blah2 blah2",
				},
			},
		})
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(string(errorContent))),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Contains(t, message, "map[a:b c:d]")
		assert.Contains(t, message, "password")
		assert.Contains(t, message, "this value should not be blank")
		assert.Contains(t, message, "email")
		assert.Contains(t, message, "blah2 blah2")
	})
}

func TestSetAttributeErrorsFromServerResponse(t *testing.T) {
	t.Run(
		"returns expected path if there are no children",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

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

			SetAttributeErrorsFromServerResponse(
				"summary",
				&httpResponse,
				&diags,
			)

			attributePath := path.Root("attribute")

			want := diag.Diagnostics{}
			want.AddAttributeError(attributePath, "summary", "error1")
			want.AddAttributeError(attributePath, "summary", "error2")

			assert.Equal(t, want, diags.Errors())
		},
	)

	t.Run(
		"returns expected path if there are children",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

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
    "attributeId": ["error1", "error2"]
  }
}
          `),
					),
				),
			}

			SetAttributeErrorsFromServerResponse(
				"summary",
				&httpResponse,
				&diags,
			)

			attributePath := path.Root("attribute").AtMapKey("id")

			want := diag.Diagnostics{}
			want.AddAttributeError(attributePath, "summary", "error1")
			want.AddAttributeError(attributePath, "summary", "error2")

			assert.Equal(t, want, diags.Errors())
		},
	)

	t.Run(
		"sets no errors if errorResponse cannot be translated",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			httpResponse := http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewReader([]byte(``))),
			}

			SetAttributeErrorsFromServerResponse(
				"summary",
				&httpResponse,
				&diags,
			)

			assert.False(t, diags.HasError())
		},
	)

	t.Run("sets no errors if httpResponse is nil", func(t *testing.T) {
		diags := diag.Diagnostics{}

		SetAttributeErrorsFromServerResponse(
			"summary",
			nil,
			&diags,
		)

		assert.False(t, diags.HasError())
	})
}

func ExampleSetAttributeErrorsFromServerResponse() {
	diags := diag.Diagnostics{}

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

	SetAttributeErrorsFromServerResponse("summary", &httpResponse, &diags)

	fmt.Println(diags.Errors())
	// Output: [{{error1 summary} {[attribute]}} {{error2 summary} {[attribute]}}]
}

func ExampleSetAttributeErrorsFromServerResponse_nested() {
	diags := diag.Diagnostics{}

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
    "attributeId": ["error1", "error2"]
  }
}
          `),
			),
		),
	}

	SetAttributeErrorsFromServerResponse("summary", &httpResponse, &diags)

	fmt.Println(diags.Errors())
	// Output: [{{error1 summary} {[attribute id]}} {{error2 summary} {[attribute id]}}]
}

func Test_normalizeErrorResponseKey(t *testing.T) {
	t.Run("camel case is normalize correctly", func(t *testing.T) {
		want := "instance_id"
		got := normalizeErrorResponseKey("instanceId")

		assert.Equal(t, want, got)
	})

	t.Run("keys with dots are normalize correctly", func(t *testing.T) {
		want := "instance_id"
		got := normalizeErrorResponseKey("instance.Id")

		assert.Equal(t, want, got)
	})

}
