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

func TestSdkError(t *testing.T) {
	t.Run("adds generic error when err is nil", func(t *testing.T) {
		diags := diag.Diagnostics{}
		httpResponse := http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}

		SdkError(context.TODO(), &diags, nil, &httpResponse)

		assert.Len(t, diags.Errors(), 1)
		assert.Equal(t, sdkErrSummary, diags.Errors()[0].Summary())
		assert.Equal(
			t,
			"An error has occurred in the program. Please consider opening an issue.",
			diags.Errors()[0].Detail(),
		)
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

		SdkError(
			context.TODO(),
			&diags,
			errors.New("error content"),
			&httpResponse,
		)

		attributePath := path.Root("name")
		want := diag.Diagnostics{}
		want.AddAttributeError(
			attributePath,
			sdkErrSummary,
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

		SdkError(
			context.TODO(),
			&diags,
			errors.New("error content"),
			&httpResponse,
		)

		assert.Len(t, diags.Errors(), 1)
		assert.Equal(t, "Unexpected Error", diags.Errors()[0].Summary())
		assert.Equal(t, defaultErrMsg, diags.Errors()[0].Detail())
	})

	t.Run("sets default error if error is nil", func(t *testing.T) {
		diags := diag.Diagnostics{}

		SdkError(context.TODO(), &diags, nil, nil)

		want := diag.Diagnostics{}
		want.AddError(sdkErrSummary, defaultErrMsg)

		assert.Equal(t, want, diags)
	})

	t.Run(
		"sets  expected error if server returns a 504 response",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			SdkError(
				context.TODO(),
				&diags,
				errors.New(""),
				&http.Response{
					Body:       io.NopCloser(bytes.NewReader([]byte("tralala"))),
					StatusCode: 504,
				},
			)

			want := diag.Diagnostics{}
			want.AddError(
				sdkErrSummary,
				"The server took too long to respond.",
			)

			assert.Equal(t, want, diags)
		},
	)

	t.Run(
		"sets error if response body cannot be mapped to errorResponse",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			SdkError(
				context.TODO(),
				&diags,
				errors.New(""),
				&http.Response{
					Body: io.NopCloser(bytes.NewReader([]byte(""))),
				},
			)

			want := diag.Diagnostics{}
			want.AddError(sdkErrSummary, defaultErrMsg)

			assert.Equal(t, want, diags)
		},
	)

	t.Run(
		"sets expected error if server returns a 404 response",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			SdkError(
				context.TODO(),
				&diags,
				errors.New(""),
				&http.Response{
					Body:       io.NopCloser(bytes.NewReader([]byte(``))),
					StatusCode: 404,
				},
			)

			want := diag.Diagnostics{}
			want.AddError(sdkErrSummary, "Resource not found.")

			assert.Equal(t, want, diags)
		},
	)

	t.Run("sets expected path if there are no children", func(t *testing.T) {
		diags := diag.Diagnostics{}

		SdkError(
			context.TODO(),
			&diags,
			errors.New(""),
			&http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte(`
						{
		              		"correlationId": "correlationId",
		              		"errorCode": "errorCode",
		              		"errorMessage": "errorMessage",
							"errorDetails":  {
								"attribute": ["error1", "error2"]
							}
		            	}`,
				))),
			},
		)

		attributePath := path.Root("attribute")
		want := diag.Diagnostics{}
		want.AddAttributeError(attributePath, sdkErrSummary, "error1")
		want.AddAttributeError(attributePath, sdkErrSummary, "error2")
		assert.Equal(t, want, diags.Errors())
	})

	t.Run("sets expected path if there are children", func(t *testing.T) {
		diags := diag.Diagnostics{}

		SdkError(
			context.TODO(),
			&diags,
			errors.New(""),
			&http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte(`
						{
		              		"correlationId": "correlationId",
		              		"errorCode": "errorCode",
		              		"errorMessage": "errorMessage",
							"errorDetails":  {
								"attributeId": ["error1", "error2"]
							}
		            	}`,
				))),
			},
		)

		attributePath := path.Root("attribute").AtMapKey("id")
		want := diag.Diagnostics{}
		want.AddAttributeError(attributePath, sdkErrSummary, "error1")
		want.AddAttributeError(attributePath, sdkErrSummary, "error2")
		assert.Equal(t, want, diags.Errors())
	})

	t.Run("camelcase errorDetails key is normalized correctly", func(t *testing.T) {
		diags := diag.Diagnostics{}

		SdkError(
			context.TODO(),
			&diags,
			errors.New(""),
			&http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte(`
						{
		              		"correlationId": "correlationId",
		              		"errorCode": "errorCode",
		              		"errorMessage": "errorMessage",
							"errorDetails":  {
								"attributeId": ["error"]
							}
		            	}`,
				))),
			},
		)

		attributePath := path.Root("attribute").AtMapKey("id")
		want := diag.Diagnostics{}
		want.AddAttributeError(attributePath, sdkErrSummary, "error")
		assert.Equal(t, want, diags.Errors())
	})

	t.Run("errorDetails with a dot is normalized correctly", func(t *testing.T) {
		diags := diag.Diagnostics{}

		SdkError(
			context.TODO(),
			&diags,
			errors.New(""),
			&http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte(`
						{
		              		"correlationId": "correlationId",
		              		"errorCode": "errorCode",
		              		"errorMessage": "errorMessage",
							"errorDetails":  {
								"attribute.Id": ["error"]
							}
		            	}`,
				))),
			},
		)

		attributePath := path.Root("attribute").AtMapKey("id")
		want := diag.Diagnostics{}
		want.AddAttributeError(attributePath, sdkErrSummary, "error")
		assert.Equal(t, want, diags.Errors())
	})

	t.Run("writes error is response is nil", func(t *testing.T) {
		diags := diag.Diagnostics{}

		SdkError(
			context.TODO(),
			&diags,
			errors.New("tralala"),
			nil,
		)

		want := diag.Diagnostics{}
		want.AddError(sdkErrSummary, "tralala")
		assert.Equal(t, want, diags)
	})

	t.Run("can handle nested errorDetails", func(t *testing.T) {
		diags := diag.Diagnostics{}

		SdkError(
			context.TODO(),
			&diags,
			errors.New(""),
			&http.Response{
				Body: io.NopCloser(
					bytes.NewReader(
						[]byte(`
						{
		              		"correlationId": "correlationId",
		              		"errorCode": "errorCode",
		              		"errorMessage": "errorMessage",
							"errorDetails":  {
								"attribute": {
									"0": ["error"]
								}
							}
		            	}`,
						),
					),
				),
			},
		)

		attributePath := path.Root("attribute")
		want := diag.Diagnostics{}
		want.AddAttributeError(attributePath, sdkErrSummary, "error")
		assert.Equal(t, want, diags.Errors())
	})

	t.Run(
		"error is set if errorDetails is not one of the expected values",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			SdkError(
				context.TODO(),
				&diags,
				errors.New(""),
				&http.Response{
					Body: io.NopCloser(
						bytes.NewReader(
							[]byte(`
						{
		              		"correlationId": "correlationId",
		              		"errorCode": "errorCode",
		              		"errorMessage": "errorMessage",
							"errorDetails":  {
								"attribute": 26
							}
		            	}`,
							),
						),
					),
				},
			)

			want := diag.Diagnostics{}
			want.AddError("Unexpected Error", defaultErrMsg)
			assert.Equal(t, want, diags.Errors())
		},
	)

	t.Run(
		"error is set if string array is expected but something else is passed",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			SdkError(
				context.TODO(),
				&diags,
				errors.New(""),
				&http.Response{
					Body: io.NopCloser(
						bytes.NewReader(
							[]byte(`
						{
		              		"correlationId": "correlationId",
		              		"errorCode": "errorCode",
		              		"errorMessage": "errorMessage",
							"errorDetails":  {
								"attribute": [26]
							}
		            	}`,
							),
						),
					),
				},
			)

			want := diag.Diagnostics{}
			want.AddError("Unexpected Error", defaultErrMsg)
			assert.Equal(t, want, diags.Errors())
		},
	)

	t.Run(
		"error is set if map of string arrays is expected map of non strings is passed",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			SdkError(
				context.TODO(),
				&diags,
				errors.New(""),
				&http.Response{
					Body: io.NopCloser(
						bytes.NewReader(
							[]byte(`
						{
		              		"correlationId": "correlationId",
		              		"errorCode": "errorCode",
		              		"errorMessage": "errorMessage",
							"errorDetails":  {
								"attribute": {
									"0": [26]
								}
							}
		            	}`,
							),
						),
					),
				},
			)

			want := diag.Diagnostics{}
			want.AddError("Unexpected Error", defaultErrMsg)
			assert.Equal(t, want, diags.Errors())
		},
	)

	t.Run(
		"error is set if map of string arrays is expected but a map of non arrays is passed",
		func(t *testing.T) {
			diags := diag.Diagnostics{}

			SdkError(
				context.TODO(),
				&diags,
				errors.New(""),
				&http.Response{
					Body: io.NopCloser(
						bytes.NewReader(
							[]byte(`
						{
		              		"correlationId": "correlationId",
		              		"errorCode": "errorCode",
		              		"errorMessage": "errorMessage",
							"errorDetails":  {
								"attribute": {
									"0": 26
								}
							}
		            	}`,
							),
						),
					),
				},
			)

			want := diag.Diagnostics{}
			want.AddError("Unexpected Error", defaultErrMsg)
			assert.Equal(t, want, diags.Errors())
		},
	)
}

func ExampleSdkError() {
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

	SdkError(context.TODO(), &diags, errors.New("error content"), &httpResponse)

	fmt.Println(diags.Errors())
	// Output: [{{the name is invalid Unexpected API Error} {[name]}}]
}

func TestGeneralError(t *testing.T) {
	diags := diag.Diagnostics{}
	GeneralError(&diags, context.TODO(), errors.New("tralala"))

	assert.Len(t, diags.Errors(), 1)
	assert.Equal(t, "Unexpected Error", diags.Errors()[0].Summary())
	assert.Equal(t, defaultErrMsg, diags.Errors()[0].Detail())
}

func ExampleGeneralError() {
	diags := diag.Diagnostics{}
	GeneralError(&diags, context.TODO(), errors.New("error content"))

	fmt.Println(diags.Errors())
	// Output: [{An error has occurred in the program. Please consider opening an issue. Unexpected Error}]
}

func TestImportOnlyError(t *testing.T) {
	diags := diag.Diagnostics{}
	ImportOnlyError(&diags)

	assert.Len(t, diags.Errors(), 1)
	assert.Equal(
		t,
		"Resource can only be imported, not created.",
		diags.Errors()[0].Summary(),
	)
	assert.Equal(t, "", diags.Errors()[0].Detail())
}

func ExampleImportOnlyError() {
	diags := diag.Diagnostics{}
	ImportOnlyError(&diags)

	fmt.Println(diags.Errors())
	// Output: [{ Resource can only be imported, not created.}]
}

func TestUnexpectedImportIdentifierError(t *testing.T) {
	diags := diag.Diagnostics{}
	UnexpectedImportIdentifierError(&diags, "format", "got")

	assert.Len(t, diags.Errors(), 1)
	assert.Equal(t, "Unexpected Import Identifier", diags.Errors()[0].Summary())
	assert.Equal(
		t,
		`Expected import identifier with format: "format". Got: "got"`,
		diags.Errors()[0].Detail(),
	)
}

func ExampleUnexpectedImportIdentifierError() {
	diags := diag.Diagnostics{}
	UnexpectedImportIdentifierError(
		&diags,
		"load_balancer_id,listener_id",
		"f6d09965-c857-4d9b-a17f-c21bf13ddcd4",
	)

	fmt.Println(diags.Errors())
	// Output: [{Expected import identifier with format: "load_balancer_id,listener_id". Got: "f6d09965-c857-4d9b-a17f-c21bf13ddcd4" Unexpected Import Identifier}]
}

func Test_writeSDKOutput(t *testing.T) {
	diags := diag.Diagnostics{}

	writeSDKOutput("tralala", &diags)

	assert.Len(t, diags.Errors(), 1)
	assert.Equal(t, sdkErrSummary, diags.Errors()[0].Summary())
	assert.Equal(t, "tralala", diags.Errors()[0].Detail())
}

func Test_handleStringErrorCollection(t *testing.T) {
	t.Run("errors are parsed correctly", func(t *testing.T) {
		diags := diag.Diagnostics{}

		got := handleStringErrorCollection(&diags, path.Root("attribute"), []interface{}{"error"})

		assert.False(t, got)
		assert.Len(t, diags.Errors(), 1)
		assert.Equal(t, sdkErrSummary, diags.Errors()[0].Summary())
		assert.Equal(t, "error", diags.Errors()[0].Detail())
	})

	t.Run("empty errors are parsed correctly", func(t *testing.T) {
		diags := diag.Diagnostics{}

		got := handleStringErrorCollection(&diags, path.Root("attribute"), []interface{}{})

		assert.False(t, got)
		assert.False(t, diags.HasError())
	})

	t.Run("returns true if error cannot be parsed", func(t *testing.T) {
		diags := diag.Diagnostics{}

		got := handleStringErrorCollection(&diags, path.Root("attribute"), []interface{}{1})

		assert.True(t, got)
		assert.False(t, diags.HasError())
	})
}
