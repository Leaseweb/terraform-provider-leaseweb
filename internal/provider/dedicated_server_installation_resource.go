package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedServer"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

var (
	_ resource.Resource              = &dedicatedServerInstallationResource{}
	_ resource.ResourceWithConfigure = &dedicatedServerInstallationResource{}
)

func NewDedicatedServerInstallationResource() resource.Resource {
	return &dedicatedServerInstallationResource{}
}

type dedicatedServerInstallationResource struct {
	// TODO: Refactor this part, apiKey shouldn't be here.
	apiKey string
	client dedicatedServer.DedicatedServerAPI
}

type dedicatedServerInstallationResourceModel struct {
	ID                types.String   `tfsdk:"id"`
	DedicatedServerID types.String   `tfsdk:"dedicated_server_id"`
	CallbackURL       types.String   `tfsdk:"callback_url"`
	ControlPanelID    types.String   `tfsdk:"control_panel_id"`
	Device            types.String   `tfsdk:"device"`
	Hostname          types.String   `tfsdk:"hostname"`
	OperatingSystemID types.String   `tfsdk:"operating_system_id"`
	Partitions        types.List     `tfsdk:"partitions"`
	Password          types.String   `tfsdk:"password"`
	PostInstallScript types.String   `tfsdk:"post_install_script"`
	PowerCycle        types.Bool     `tfsdk:"power_cycle"`
	Raid              types.Object   `tfsdk:"raid"`
	SSHKeys           []types.String `tfsdk:"ssh_keys"`
	Timezone          types.String   `tfsdk:"timezone"`
}

type dedicatedServerInstallationRaidModel struct {
	Level         types.Int32  `tfsdk:"level"`
	NumberOfDisks types.Int32  `tfsdk:"number_of_disks"`
	Type          types.String `tfsdk:"type"`
}

type dedicatedServerInstallationPartitionsModel struct {
	Filesystem types.String `tfsdk:"filesystem"`
	Mountpoint types.String `tfsdk:"mountpoint"`
	Size       types.String `tfsdk:"size"`
}

func (m dedicatedServerInstallationPartitionsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"filesystem": types.StringType,
		"mountpoint": types.StringType,
		"size":       types.StringType,
	}
}

func (r *dedicatedServerInstallationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_server_installation"
}

func (d *dedicatedServerInstallationResource) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		dedicatedServer.ContextAPIKeys,
		map[string]dedicatedServer.APIKey{
			"X-LSW-Auth": {Key: d.apiKey, Prefix: ""},
		},
	)
}

func (d *dedicatedServerInstallationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	configuration := dedicatedServer.NewConfiguration()

	// TODO: Refactor this part, ProviderData can be managed directly, not within client.
	coreClient, ok := req.ProviderData.(client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return
	}
	d.apiKey = coreClient.ProviderData.ApiKey
	if coreClient.ProviderData.Host != nil {
		configuration.Host = *coreClient.ProviderData.Host
	}
	if coreClient.ProviderData.Scheme != nil {
		configuration.Scheme = *coreClient.ProviderData.Scheme
	}

	apiClient := dedicatedServer.NewAPIClient(configuration)
	d.client = apiClient.DedicatedServerAPI
}

func (r *dedicatedServerInstallationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	raid := func() schema.SingleNestedAttribute {
		return schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"level": schema.Int32Attribute{
					Description: "RAID level to apply to your installation, this value is only required if you specify a type HW or SW. Valid options are \n  - *0*\n  - *1*\n  - *5*\n  - *10*\n",
					Optional:    true,
					Validators: []validator.Int32{
						int32validator.OneOf([]int32{0, 1, 5, 10}...),
					},
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.RequiresReplace(),
					},
				},
				"number_of_disks": schema.Int32Attribute{
					Description: "The number of disks you want to apply RAID on. If not specified all disks are used",
					Optional:    true,
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.RequiresReplace(),
					},
				},
				"type": schema.StringAttribute{
					Description: "RAID type to apply to your installation. NONE is the equivalent of pass-through mode on HW RAID equipped servers. Valid options are \n  - *HW*\n  - *SW*\n  - *NONE*\n",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf([]string{"HW", "SW", "NONE"}...),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
		}
	}

	partitions := func() schema.ListNestedAttribute {
		return schema.ListNestedAttribute{
			Optional: true,
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"filesystem": schema.StringAttribute{
						Description: "File system in which partition would be mounted",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"mountpoint": schema.StringAttribute{
						Description: "The partition mount point (eg /home). Mandatory for the root partition (/) and not intended to be used in swap partition",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"size": schema.StringAttribute{
						Description: "Size of the partition (Normally in MB, but this is OS-specific)",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
		}
	}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the installation job",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dedicated_server_id": schema.StringAttribute{
				Description: "The ID of a server",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"callback_url": schema.StringAttribute{
				Description: "Url which will receive callbacks when the installation is finished or failed",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"control_panel_id": schema.StringAttribute{
				Description: "Control panel identifier",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device": schema.StringAttribute{
				Description: `Block devices in a disk set in which the partitions will be installed. Supported values are any disk set id, ` + "`SATA_SAS`" + ` or ` + "`NVME`" + `.`,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname to be used in your installation",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"operating_system_id": schema.StringAttribute{
				Description: "Operating system identifier",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"partitions": partitions(),
			"password": schema.StringAttribute{
				Description: "Server root password. If not provided, it would be automatically generated",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"post_install_script": schema.StringAttribute{
				Description: "Base64 Encoded string containing a valid bash script to be run right after the installation",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"power_cycle": schema.BoolAttribute{
				Description: "If true, allows system reboots to happen automatically within the process. Otherwise, you should do them manually",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"raid": raid(),
			"ssh_keys": schema.SetAttribute{
				Description: "List of public sshKeys to be setup in your installation",
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"timezone": schema.StringAttribute{
				Description: "Timezone represented as Geographical_Area/City",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}

}

func (r *dedicatedServerInstallationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dedicatedServerInstallationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	// Extract the Raid configuration from the plan
	var raidPlan dedicatedServerInstallationRaidModel
	plan.Raid.As(ctx, &raidPlan, basetypes.ObjectAsOptions{})

	// Extract Partitions configuration from the plan into a Go slice
	partitionsPlan := make([]dedicatedServerInstallationPartitionsModel, 0, len(plan.Partitions.Elements()))
	plan.Partitions.ElementsAs(ctx, &partitionsPlan, false)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preparing partitions for the installation options.
	var partitions []dedicatedServer.Partition
	if !plan.Partitions.IsNull() && !plan.Partitions.IsUnknown() {

		for _, p := range partitionsPlan {
			// Check if all fields are either null or unknown, if so, skip
			if valueStringOrNil(p.Mountpoint) == nil &&
				valueStringOrNil(p.Size) == nil &&
				valueStringOrNil(p.Filesystem) == nil {
				continue
			}

			partitions = append(partitions, dedicatedServer.Partition{
				Filesystem: valueStringOrNil(p.Filesystem),
				Size:       valueStringOrNil(p.Size),
				Mountpoint: valueStringOrNil(p.Mountpoint),
			})
		}

	}

	// Preparing RAID configuration for the installation options
	var raid *dedicatedServer.Raid
	// Check that at least one RAID field is set before initializing the RAID struct.
	if !plan.Raid.IsNull() && !plan.Raid.IsUnknown() &&
		(valueInt32OrNil(raidPlan.Level) != nil ||
			valueInt32OrNil(raidPlan.NumberOfDisks) != nil ||
			valueStringOrNil(raidPlan.Type) != nil) {

		raid = &dedicatedServer.Raid{
			Level:         (*dedicatedServer.RaidLevel)(valueInt32OrNil(raidPlan.Level)),
			NumberOfDisks: valueInt32OrNil(raidPlan.NumberOfDisks),
			Type:          (*dedicatedServer.RaidType)(valueStringOrNil(raidPlan.Type)),
		}
	}

	// Preparing SSH keys for the installation options, combining them into a single string
	var SSHKeysList []string
	for _, k := range plan.SSHKeys {
		if valueStringOrNil(k) != nil {
			SSHKeysList = append(SSHKeysList, k.ValueString())
		}
	}
	SSHKeys := strings.Join(SSHKeysList, "\n")

	opts := dedicatedServer.NewInstallOperatingSystemOpts(plan.OperatingSystemID.ValueString())
	opts.CallbackUrl = valueStringOrNil(plan.CallbackURL)
	opts.ControlPanelId = valueStringOrNil(plan.ControlPanelID)
	opts.Device = valueStringOrNil(plan.Device)
	opts.Hostname = valueStringOrNil(plan.Hostname)
	opts.Partitions = partitions
	opts.Password = valueStringOrNil(plan.Password)
	opts.PostInstallScript = valueStringOrNil(plan.PostInstallScript)
	opts.PowerCycle = valueBoolOrNil(plan.PowerCycle)
	opts.Raid = raid
	opts.Timezone = valueStringOrNil(plan.Timezone)
	if len(SSHKeysList) > 0 {
		opts.SshKeys = &SSHKeys
	}

	serverID := plan.DedicatedServerID.ValueString()
	result, response, err := r.client.InstallOperatingSystem(r.authContext(ctx), serverID).
		InstallOperatingSystemOpts(*opts).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error resource dedicated_server_installation for server id: %q", serverID),
			extractErrorMessage(response, err),
		)
		return
	}

	payload := result.GetPayload()
	plan.ID = types.StringValue(result.GetUuid())
	plan.Device = types.StringValue(payload.GetDevice())
	plan.Timezone = types.StringValue(payload.GetTimezone())
	plan.PowerCycle = types.BoolValue(payload.GetPowerCycle())

	// Preparing and converting partitions into types.Object to store in the state
	var partitionsObjects []attr.Value
	for _, p := range payload.GetPartitions() {
		partition := dedicatedServerInstallationPartitionsModel{
			Filesystem: types.StringValue(p.GetFilesystem()),
			Mountpoint: types.StringValue(p.GetMountpoint()),
			Size:       types.StringValue(p.GetSize()),
		}

		partitionObj, diags := types.ObjectValueFrom(ctx, partition.AttributeTypes(), partition)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		partitionsObjects = append(partitionsObjects, partitionObj)
	}

	// Convert the slice of partition objects to a types.List and store it in the plan
	partitionsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: dedicatedServerInstallationPartitionsModel{}.AttributeTypes()}, partitionsObjects)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	plan.Partitions = partitionsList

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

}

func (r *dedicatedServerInstallationResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *dedicatedServerInstallationResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *dedicatedServerInstallationResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// TODO: Should goes like helper/shared.
func valueStringOrNil(s basetypes.StringValue) *string {
	if s.IsUnknown() {
		return nil
	}
	return s.ValueStringPointer()
}

func valueBoolOrNil(s basetypes.BoolValue) *bool {
	if s.IsUnknown() {
		return nil
	}
	return s.ValueBoolPointer()
}

func valueInt32OrNil(s basetypes.Int32Value) *int32 {
	if s.IsUnknown() {
		return nil
	}
	return s.ValueInt32Pointer()
}

func extractErrorMessage(resp *http.Response, err error) string {
	if resp == nil || resp.Body == nil || resp.StatusCode < 400 {
		return err.Error()
	}

	defer resp.Body.Close() // Ensure the body is closed
	var errorResponse map[string]interface{}

	// Attempt to decode the response body as JSON
	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
		if errorMessage, ok := errorResponse["errorMessage"]; ok {
			return fmt.Sprintf("%v", errorMessage)
		}
	}

	return err.Error()
}
