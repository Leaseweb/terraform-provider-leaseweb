package logging

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/repository"
	"github.com/stretchr/testify/assert"
)

func TestServiceError(t *testing.T) {
	t.Run("response is set", func(t *testing.T) {

		diags := diag.Diagnostics{}

		ServiceError(
			context.TODO(),
			&repository.ErrorResponse{},
			&diags,
			"summary",
			"detail",
		)

		assert.Equal(
			t,
			"summary",
			diags[0].Summary(),
			"error contains summary",
		)
		assert.Equal(
			t,
			"detail",
			diags[0].Detail(),
			"error contains detail",
		)
	})

	t.Run("response is not set", func(t *testing.T) {
		diags := diag.Diagnostics{}
		summary := "summary"
		detail := "detail"

		ServiceError(context.TODO(), nil, &diags, summary, detail)

		assert.Equal(
			t,
			"summary",
			diags[0].Summary(),
			"error contains summary",
		)
		assert.Equal(
			t,
			"detail",
			diags[0].Detail(),
			"error contains detail",
		)
	})
}
