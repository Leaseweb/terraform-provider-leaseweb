// Package logging implements logging that the end user sees.
package logging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud/shared/repository"
)

// FacadeError prints the passed errorResponse as a Terraform error log.
func FacadeError(
	ctx context.Context,
	errorResponse *repository.ErrorResponse,
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
