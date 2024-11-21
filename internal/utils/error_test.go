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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Run("adds generic error when err is nil", func(t *testing.T) {
		diags := diag.Diagnostics{}
		httpResponse := http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}

		Error(context.TODO(), &diags, "summary", nil, &httpResponse)

		assert.Len(t, diags.Errors(), 1)
		assert.Equal(t, "summary", diags.Errors()[0].Summary())
		assert.Equal(
			t,
			"An error has occurred in the program. Please consider opening an issue.",
			diags.Errors()[0].Detail(),
		)
	})

	t.Run("handles SDK error when HTTP status is non-error", func(t *testing.T) {
		diags := diag.Diagnostics{}
		httpResponse := http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"key":"value"}`))),
		}

		Error(
			context.TODO(),
			&diags,
			"summary",
			errors.New("SDK error: enum doesn't match"),
			&httpResponse,
		)

		assert.Len(t, diags.Errors(), 1)
		assert.Equal(t, "summary", diags.Errors()[0].Summary())
		assert.Equal(t, "SDK error: enum doesn't match", diags.Errors()[0].Detail())
	})

	t.Run("adds attribute error from HTTP validation details", func(t *testing.T) {
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
                "name": ["the name is invalid"]
              }
            }
          `),
				),
			),
		}

		Error(
			context.TODO(),
			&diags,
			"summary",
			errors.New("error content"),
			&httpResponse,
		)

		attributePath := path.Root("name")
		want := diag.Diagnostics{}
		want.AddAttributeError(
			attributePath,
			"summary",
			"the name is invalid",
		)
		assert.Equal(t, want, diags.Errors())
	})

	t.Run("handles regular HTTP error response", func(t *testing.T) {
		diags := diag.Diagnostics{}
		httpResponse := http.Response{
			StatusCode: 500,
			Body: io.NopCloser(
				bytes.NewReader(
					[]byte(`
            {
              "correlationId": "correlationId",
              "errorCode": "404",
              "errorMessage": "Server not found"
            }
          `),
				),
			),
		}

		Error(
			context.TODO(),
			&diags,
			"summary",
			errors.New("error content"),
			&httpResponse,
		)

		assert.Len(t, diags.Errors(), 1)
		assert.Equal(t, "summary", diags.Errors()[0].Summary())
		assert.Equal(
			t,
			"{\n  \"errorCode\": \"404\",\n  \"errorMessage\": \"Server not found\",\n  \"correlationId\": \"correlationId\"\n}",
			diags.Errors()[0].Detail(),
		)
	})
}

func ExampleError() {
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
  "name": ["the name is invalid"]
  }
}
          `),
			),
		),
	}

	Error(
		context.TODO(),
		&diags,
		"summary",
		errors.New("error content"),
		&httpResponse,
	)

	fmt.Println(diags.Errors())
	// Output: [{{the name is invalid summary} {[name]}}]
}

func Test_mapErrorDetailsKey(t *testing.T) {
	t.Run("camel case is normalize correctly", func(t *testing.T) {
		want := "instance_id"
		got := mapErrorDetailsKey("instanceId")
		assert.Equal(t, want, got)
	})
	t.Run("keys with dots are normalize correctly", func(t *testing.T) {
		want := "instance_id"
		got := mapErrorDetailsKey("instance.Id")
		assert.Equal(t, want, got)
	})
}

func Test_errorHandler_handleValidationErrors(t *testing.T) {
	t.Run("sets expected path if there are no children", func(t *testing.T) {
		diags := diag.Diagnostics{}

		errorResponse := errorResponse{
			ErrorDetails: map[string][]string{
				"attribute": {"error1", "error2"},
			},
		}
		errorHandler := errorHandler{
			summary: "summary",
			diags:   &diags,
		}
		errorHandler.handleValidationErrors(errorResponse)

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

			errorResponse := errorResponse{
				ErrorDetails: map[string][]string{
					"attributeId": {"error1", "error2"},
				},
			}
			errorHandler := errorHandler{
				summary: "summary",
				diags:   &diags,
			}
			errorHandler.handleValidationErrors(errorResponse)

			attributePath := path.Root("attribute").AtMapKey("id")
			want := diag.Diagnostics{}
			want.AddAttributeError(attributePath, "summary", "error1")
			want.AddAttributeError(attributePath, "summary", "error2")
			assert.Equal(t, want, diags.Errors())
		},
	)

	/**
	  t.Run(
	  	"sets expected error if errorResponse cannot be translated",
	  	func(t *testing.T) {
	  		diags := diag.Diagnostics{}
	  		httpResponse := http.Response{
	  			StatusCode: 500,
	  			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
	  		}
	  		Error(
	  			context.TODO(),
	  			&diags,
	  			"summary",
	  			errors.New(""),
	  			&httpResponse,
	  		)
	  		assert.Len(t, diags.Errors(), 1)
	  		assert.Equal(t, "summary", diags.Errors()[0].Summary())
	  		assert.Equal(
	  			t,
	  			DefaultErrMsg,
	  			diags.Errors()[0].Detail(),
	  		)
	  	},
	  )
	*/

	/**
	  t.Run("error is outputted if httpResponse is empty", func(t *testing.T) {
	  	diags := diag.Diagnostics{}
	  	Error(
	  		context.TODO(),
	  		&diags,
	  		"summary",
	  		errors.New("tralala"),
	  		nil,
	  	)
	  	assert.Len(t, diags.Errors(), 1)
	  	assert.Equal(t, "summary", diags.Errors()[0].Summary())
	  	assert.Equal(t, "tralala", diags.Errors()[0].Detail())
	  })
	*/

	/**
	  t.Run(
	  	"global error is set if no attribute errors can be set",
	  	func(t *testing.T) {
	  		diags := diag.Diagnostics{}
	  		httpBody := `{
	     "correlationId": "correlationId",
	     "errorCode": "errorCode",
	     "errorMessage": "errorMessage"
	    }`
	  		httpResponse := http.Response{
	  			StatusCode: 500,
	  			Body:       io.NopCloser(bytes.NewReader([]byte(httpBody))),
	  		}

	  		Error(
	  			context.TODO(),
	  			&diags,
	  			"summary",
	  			nil,
	  			&httpResponse,
	  		)

	  		assert.Len(t, diags.Errors(), 1)
	  		assert.Equal(t, "summary", diags.Errors()[0].Summary())
	  		assert.Equal(t, DefaultErrMsg, diags.Errors()[0].Detail())
	  	},
	  )
	*/
}

func Test_errorHandler_processErrorResponse(t *testing.T) {
	t.Run(
		"does not set default error if errorResponse has errors in it",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			errorResponse := errorResponse{
				ErrorDetails: map[string][]string{
					"attribute": {"error"},
				},
			}
			errorHandler := errorHandler{
				summary: "summary",
				diags:   &diags,
			}
			errorHandler.processErrorResponse(errorResponse)

			attributePath := path.Root("attribute")
			want := diag.Diagnostics{}
			want.AddAttributeError(attributePath, "summary", "error")

			assert.Equal(t, want, diags)
		},
	)

	t.Run(
		"error response if outputted if it cannot be parsed",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			errorResponse := errorResponse{}
			errorHandler := errorHandler{
				summary: "summary",
				diags:   &diags,
			}
			errorHandler.processErrorResponse(errorResponse)

			want := diag.Diagnostics{}
			want.AddError("summary", "{}")

			assert.Equal(t, want, diags)
		},
	)
}

func Test_errorHandler_handleHTTPError(t *testing.T) {
	t.Run(
		"sets error if response body cannot be mapped to errorResponse",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			errorHandler := errorHandler{
				summary: "summary",
				diags:   &diags,
				resp: &http.Response{
					Body: io.NopCloser(bytes.NewReader([]byte(``))),
				},
				ctx: context.TODO(),
			}
			errorHandler.handleHTTPError()

			want := diag.Diagnostics{}
			want.AddError("summary", DefaultErrMsg)

			assert.Equal(t, want, diags)
		},
	)
}

func Test_errorHandler_report(t *testing.T) {
	t.Run("sets default error if error is nil", func(t *testing.T) {
		diags := diag.Diagnostics{}

		errorHandler := errorHandler{
			summary: "summary",
			diags:   &diags,
			ctx:     context.TODO(),
		}
		errorHandler.report()

		want := diag.Diagnostics{}
		want.AddError("summary", DefaultErrMsg)

		assert.Equal(t, want, diags)
	})

}
