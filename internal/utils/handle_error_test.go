package utils

import (
	"bytes"
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestHandleError_responseIsSet(t *testing.T) {
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
}
