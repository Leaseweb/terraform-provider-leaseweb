package publiccloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ resource.ResourceWithConfigure = &targetGroupResource{}
)

type targetGroupResource struct {
	name   string
	client publicCloud.PublicCloudAPI
}

func (t *targetGroupResource) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = fmt.Sprintf(
		"%s_%s",
		request.ProviderTypeName,
		t.name,
	)
}

func (t *targetGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Once created, an image resource cannot be deleted via Terraform",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Can be either an Operating System or a UUID in case of a Custom Image",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				Description: `
The id of the instance which the custom image is based on. The following rules apply:
  - instance exists for instanceId
  - instance has state *STOPPED*
  - instance has a maximum rootDiskSize of 100 GB
  - instance OS must not be *windows*`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Custom image name",
			},
			"custom": schema.BoolAttribute{
				Computed:    true,
				Description: "Standard or Custom image",
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"market_apps": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"storage_types": schema.ListAttribute{
				Computed:    true,
				Description: "The supported storage types for the instance type",
				ElementType: types.StringType,
			},
			"flavour": schema.StringAttribute{
				Computed: true,
			},
			"region": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (t *targetGroupResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t *targetGroupResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (t *targetGroupResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t *targetGroupResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func (t *targetGroupResource) Configure(
	_ context.Context,
	request resource.ConfigureRequest,
	response *resource.ConfigureResponse,
) {
	if request.ProviderData == nil {
		return
	}

	coreClient, ok := request.ProviderData.(client.Client)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				request.ProviderData,
			),
		)

		return
	}

	t.client = coreClient.PublicCloudAPI
}

func NewTargetGroupResource() resource.Resource {
	return &targetGroupResource{
		name: "public_cloud_target_group",
	}
}
