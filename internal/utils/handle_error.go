package utils

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func HandleError(
	ctx context.Context,
	response *http.Response,
	diags *diag.Diagnostics,
	summary string,
	detail string,
) {
	if response != nil {
		buf := new(strings.Builder)
		_, sdkResponseError := io.Copy(buf, response.Body)
		if sdkResponseError == nil {
			tflog.Debug(
				ctx,
				"API response",
				map[string]interface{}{"response": buf.String()},
			)
		}
	}

	diags.AddError(summary, detail)
}
