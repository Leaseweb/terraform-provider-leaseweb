package logging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-leaseweb/internal/shared"
)

func HandleError(
	ctx context.Context,
	errorResponse *shared.ErrorResponse,
	diags *diag.Diagnostics,
	summary string,
	detail string,
) {
	if errorResponse != nil {
		tflog.Error(
			ctx,
			summary,
			map[string]interface{}{"ErrorResponse": errorResponse},
		)
	}

	diags.AddError(summary, detail)
}
