package utils

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"strings"
)

func HandleError(
	ctx context.Context,
	response *http.Response,
	diags *diag.Diagnostics,
	summary string,
	detail string,
) {
	buf := new(strings.Builder)
	_, sdkResponseError := io.Copy(buf, response.Body)
	if sdkResponseError == nil {
		tflog.Debug(
			ctx,
			"API response",
			map[string]interface{}{"response": buf.String()},
		)
	}

	diags.AddError(summary, detail)
}
