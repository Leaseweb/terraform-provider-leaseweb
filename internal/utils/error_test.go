package utils

import (
	"bytes"
	"context"
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

		want := "{\n  \"errorMessage\": \"this is error message\"\n}"
		got := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, want, got)
	})

	t.Run("errorMessage must return if response contains object of errorMessage", func(t *testing.T) {
		errorContent, _ := json.Marshal(map[string]interface{}{"errorMessage": map[string]string{"a": "b", "c": "d"}})
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(string(errorContent))),
		}

		want := "{\n  \"errorMessage\": {\n    \"a\": \"b\",\n    \"c\": \"d\"\n  }\n}"
		got := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, want, got)
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

		want := "{\n  \"errorDetails\": {\n    \"email\": [\n      \"this value should be valid\",\n      \"blah2 blah2\"\n    ],\n    \"password\": [\n      \"this value should not be blank\",\n      \"blah blah\"\n    ]\n  },\n  \"errorMessage\": {\n    \"a\": \"b\",\n    \"c\": \"d\"\n  }\n}"
		got := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, want, got)
	})
}

func TestHandleSdkError(t *testing.T) {
	t.Run("sets expected path if there are no children", func(t *testing.T) {
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

		HandleSdkError("summary", &httpResponse, nil, &diags, context.TODO())

		attributePath := path.Root("attribute")

		want := diag.Diagnostics{}
		want.AddAttributeError(attributePath, "summary", "error1")
		want.AddAttributeError(attributePath, "summary", "error2")

		assert.Equal(t, want, diags.Errors())
	})

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

			HandleSdkError("summary", &httpResponse, nil, &diags, context.TODO())

			attributePath := path.Root("attribute").AtMapKey("id")

			want := diag.Diagnostics{}
			want.AddAttributeError(attributePath, "summary", "error1")
			want.AddAttributeError(attributePath, "summary", "error2")

			assert.Equal(t, want, diags.Errors())
		},
	)

	t.Run(
		"sets expected error if errorResponse cannot be translated",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			httpResponse := http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewReader([]byte(``))),
			}

			HandleSdkError("summary", &httpResponse, nil, &diags, context.TODO())

			assert.Len(t, diags.Errors(), 1)
			assert.Equal(t, "summary", diags.Errors()[0].Summary())
			assert.Equal(
				t,
				"unexpected end of JSON input",
				diags.Errors()[0].Detail(),
			)
		},
	)

	t.Run("error is outputted if httpResponse is empty", func(t *testing.T) {
		diags := diag.Diagnostics{}

		HandleSdkError(
			"summary",
			nil,
			errors.New("tralala"),
			&diags,
			context.TODO(),
		)

		assert.Len(t, diags.Errors(), 1)
		assert.Equal(t, "summary", diags.Errors()[0].Summary())
		assert.Equal(t, "tralala", diags.Errors()[0].Detail())
	})
}

func ExampleHandleSdkError() {
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

	HandleSdkError("summary", &httpResponse, nil, &diags, context.TODO())

	fmt.Println(diags.Errors())
	// Output: [{{error1 summary} {[attribute]}} {{error2 summary} {[attribute]}}]
}

func ExampleHandleSdkError_nested() {
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

	HandleSdkError("summary", &httpResponse, nil, &diags, context.TODO())

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

func Test_handleError(t *testing.T) {
	t.Run("expected log is written when an error is passed", func(t *testing.T) {
		diags := diag.Diagnostics{}
		err := errors.New("tralala")
		handleError("summary", err, &diags)

		want := diag.Diagnostics{}
		want.AddError("summary", "tralala")

		assert.Equal(t, want, diags.Errors())
	})

	t.Run(
		"expected log is written when an error is not passed",
		func(t *testing.T) {
			diags := diag.Diagnostics{}
			handleError("summary", nil, &diags)

			assert.Len(t, diags.Errors(), 1)
			assert.Contains(t, diags.Errors()[0].Summary(), "summary")
			assert.Contains(t, diags.Errors()[0].Detail(), DefaultErrMsg)
		},
	)
}
