package utils

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	t.Run("response is set", func(t *testing.T) {
		response := http.Response{}
		response.Body = io.NopCloser(bytes.NewReader([]byte("")))

		diags := diag.Diagnostics{}
		summary := "summary"
		detail := "detail"

		HandleError(context.TODO(), &response, &diags, summary, detail)

		assert.Equal(
			t,
			"summary",
			diags[0].Summary(),
			"summary must be set",
		)
		assert.Equal(
			t,
			"detail",
			diags[0].Detail(),
			"detail must be set",
		)
	})

	t.Run("response is not set", func(t *testing.T) {
		diags := diag.Diagnostics{}
		summary := "summary"
		detail := "detail"

		HandleError(context.TODO(), nil, &diags, summary, detail)

		assert.Equal(
			t,
			"summary",
			diags[0].Summary(),
			"summary must be set",
		)
		assert.Equal(
			t,
			"detail",
			diags[0].Detail(),
			"detail must be set",
		)
	})
}
