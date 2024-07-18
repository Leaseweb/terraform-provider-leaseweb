package logging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func HandleError(
	ctx context.Context,
	response *string,
	diags *diag.Diagnostics,
	summary string,
	detail string,
) {
	if response != nil {
		tflog.Debug(
			ctx,
			"API response",
			map[string]interface{}{"response": response},
		)
	}

	diags.AddError(summary, detail)
}
