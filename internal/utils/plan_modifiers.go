package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var (
	_ planmodifier.String = &TralalaModifier{}
)

type TralalaModifier struct{}

func (t TralalaModifier) Description(_ context.Context) string {
	return "tralala"
}

func (t TralalaModifier) MarkdownDescription(ctx context.Context) string {
	return t.Description(ctx)
}

func (t TralalaModifier) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	response.PlanValue = request.StateValue
}
