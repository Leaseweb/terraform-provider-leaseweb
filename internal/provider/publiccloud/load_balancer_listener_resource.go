package publiccloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/v2/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &loadBalancerListenerResource{}
	_ resource.ResourceWithImportState = &loadBalancerListenerResource{}
)

type loadBalancerListenerDefaultRuleResourceModel struct {
	TargetGroupID types.String `tfsdk:"target_group_id"`
}

func (l loadBalancerListenerDefaultRuleResourceModel) generateLoadBalancerListenerDefaultRule() publiccloud.LoadBalancerListenerDefaultRule {
	return *publiccloud.NewLoadBalancerListenerDefaultRule(l.TargetGroupID.ValueString())
}

func (l loadBalancerListenerDefaultRuleResourceModel) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_group_id": types.StringType,
	}
}

func adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource(loadBalancerListenerRule publiccloud.LoadBalancerListenerRule) loadBalancerListenerDefaultRuleResourceModel {
	return loadBalancerListenerDefaultRuleResourceModel{
		TargetGroupID: basetypes.NewStringValue(loadBalancerListenerRule.GetTargetGroupId()),
	}
}

type loadBalancerListenerCertificateResourceModel struct {
	PrivateKey  types.String `tfsdk:"private_key"`
	Certificate types.String `tfsdk:"certificate"`
	Chain       types.String `tfsdk:"chain"`
}

func (l loadBalancerListenerCertificateResourceModel) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"private_key": types.StringType,
		"certificate": types.StringType,
		"chain":       types.StringType,
	}
}

func (l loadBalancerListenerCertificateResourceModel) generateSslCertificate() publiccloud.SslCertificate {
	sslCertificate := publiccloud.NewSslCertificate(
		l.PrivateKey.ValueString(),
		l.Certificate.ValueString(),
	)
	if !l.Chain.IsNull() && l.Chain.ValueString() != "" {
		sslCertificate.SetChain(l.Chain.ValueString())
	}

	return *sslCertificate
}

func adaptSslCertificateToLoadBalancerListenerCertificateResource(sslCertificate publiccloud.SslCertificate) loadBalancerListenerCertificateResourceModel {
	listener := loadBalancerListenerCertificateResourceModel{
		PrivateKey:  basetypes.NewStringValue(sslCertificate.GetPrivateKey()),
		Certificate: basetypes.NewStringValue(sslCertificate.GetCertificate()),
	}

	chain, _ := sslCertificate.GetChainOk()
	if chain != nil && *chain != "" {
		listener.Chain = basetypes.NewStringPointerValue(chain)
	}

	return listener
}

type loadBalancerListenerResourceModel struct {
	ListenerID     types.String `tfsdk:"listener_id"`
	LoadBalancerID types.String `tfsdk:"load_balancer_id"`
	Protocol       types.String `tfsdk:"protocol"`
	Port           types.Int32  `tfsdk:"port"`
	Certificate    types.Object `tfsdk:"certificate"`
	DefaultRule    types.Object `tfsdk:"default_rule"`
}

func (l loadBalancerListenerResourceModel) generateLoadBalancerListenerCreateOpts(ctx context.Context) (
	*publiccloud.LoadBalancerListenerCreateOpts,
	error,
) {
	defaultRule := loadBalancerListenerDefaultRuleResourceModel{}
	defaultRuleDiags := l.DefaultRule.As(ctx, &defaultRule, basetypes.ObjectAsOptions{})
	if defaultRuleDiags != nil {
		return nil, utils.ReturnError("generateLoadBalancerListenerCreateOpts", defaultRuleDiags)
	}

	opts := publiccloud.NewLoadBalancerListenerCreateOpts(
		publiccloud.Protocol(l.Protocol.ValueString()),
		l.Port.ValueInt32(),
		defaultRule.generateLoadBalancerListenerDefaultRule(),
	)

	if !l.Certificate.IsNull() {
		certificate := loadBalancerListenerCertificateResourceModel{}
		certificateDiags := l.Certificate.As(ctx, &certificate, basetypes.ObjectAsOptions{})
		if certificateDiags != nil {
			return nil, utils.ReturnError("generateLoadBalancerListenerCreateOpts", certificateDiags)
		}

		opts.SetCertificate(certificate.generateSslCertificate())
	}

	return opts, nil
}

func (l loadBalancerListenerResourceModel) generateLoadBalancerListenerUpdateOpts(ctx context.Context) (
	*publiccloud.LoadBalancerListenerOpts,
	error,
) {
	opts := publiccloud.NewLoadBalancerListenerOpts()
	opts.SetProtocol(publiccloud.Protocol(l.Protocol.ValueString()))
	opts.SetPort(l.Port.ValueInt32())

	if !l.Certificate.IsNull() {
		certificate := loadBalancerListenerCertificateResourceModel{}
		certificateDiags := l.Certificate.As(
			ctx,
			&certificate,
			basetypes.ObjectAsOptions{},
		)
		if certificateDiags != nil {
			return nil, utils.ReturnError(
				"generateLoadBalancerListenerUpdateOpts",
				certificateDiags,
			)
		}

		opts.SetCertificate(certificate.generateSslCertificate())
	}

	if !l.DefaultRule.IsNull() {
		defaultRule := loadBalancerListenerDefaultRuleResourceModel{}
		defaultRuleDiags := l.DefaultRule.As(
			ctx,
			&defaultRule,
			basetypes.ObjectAsOptions{},
		)
		if defaultRuleDiags != nil {
			return nil, utils.ReturnError(
				"generateLoadBalancerListenerUpdateOpts",
				defaultRuleDiags,
			)
		}

		opts.SetDefaultRule(defaultRule.generateLoadBalancerListenerDefaultRule())
	}

	return opts, nil
}

func adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource(
	loadBalancerListenerDetails publiccloud.LoadBalancerListenerDetails,
	ctx context.Context,
) (*loadBalancerListenerResourceModel, error) {
	listener := loadBalancerListenerResourceModel{
		ListenerID: basetypes.NewStringValue(loadBalancerListenerDetails.GetId()),
		Protocol:   basetypes.NewStringValue(string(loadBalancerListenerDetails.GetProtocol())),
		Port:       basetypes.NewInt32Value(loadBalancerListenerDetails.GetPort()),
	}

	if len(loadBalancerListenerDetails.SslCertificates) > 0 {
		certificate, err := utils.AdaptSdkModelToResourceObject(
			loadBalancerListenerDetails.SslCertificates[0],
			loadBalancerListenerCertificateResourceModel{}.attributeTypes(),
			ctx,
			adaptSslCertificateToLoadBalancerListenerCertificateResource,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource: %w",
				err,
			)
		}
		listener.Certificate = certificate
	}

	if len(loadBalancerListenerDetails.Rules) > 0 {
		defaultRule, err := utils.AdaptSdkModelToResourceObject(
			loadBalancerListenerDetails.Rules[0],
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			ctx,
			adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource: %w",
				err,
			)
		}
		listener.DefaultRule = defaultRule
	}

	return &listener, nil
}

func adaptLoadBalancerListenerToLoadBalancerListenerResource(
	loadBalancerListener publiccloud.LoadBalancerListener,
	ctx context.Context,
) (*loadBalancerListenerResourceModel, error) {
	listener := loadBalancerListenerResourceModel{
		ListenerID: basetypes.NewStringValue(loadBalancerListener.GetId()),
		Protocol:   basetypes.NewStringValue(string(loadBalancerListener.Protocol)),
		Port:       basetypes.NewInt32Value(loadBalancerListener.GetPort()),
	}

	if len(loadBalancerListener.Rules) > 0 {
		defaultRule, err := utils.AdaptSdkModelToResourceObject(
			loadBalancerListener.Rules[0],
			loadBalancerListenerDefaultRuleResourceModel{}.attributeTypes(),
			ctx,
			adaptLoadBalancerListenerRuleToLoadBalancerListenerDefaultRuleResource,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"adaptLoadBalancerListenerToLoadBalancerListenerResource: %w",
				err,
			)
		}
		listener.DefaultRule = defaultRule
	}

	return &listener, nil
}

type loadBalancerListenerResource struct {
	utils.PubliccloudResourceAPI
}

func (l *loadBalancerListenerResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Description: utils.BetaDescription,
		Attributes: map[string]schema.Attribute{
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer ID",
			},
			"listener_id": schema.StringAttribute{
				Computed:    true,
				Description: "Listener ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Valid options are " + utils.StringTypeArrayToMarkdown(publiccloud.AllowedProtocolEnumValues),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(publiccloud.AllowedProtocolEnumValues)...),
				},
			},
			"port": schema.Int32Attribute{
				Required:    true,
				Description: "Port that the listener listens to",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"certificate": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Required only if protocol is HTTPS",
				Attributes: map[string]schema.Attribute{
					"private_key": schema.StringAttribute{
						Optional:    true,
						Description: "Client Private Key. Required only if protocol is `HTTPS`",
						Sensitive:   true,
					},
					"certificate": schema.StringAttribute{
						Optional:    true,
						Description: "Client Certificate. Required only if protocol is `HTTPS`",
						Sensitive:   true,
					},
					"chain": schema.StringAttribute{
						Optional:    true,
						Description: "CA certificate. Not required, but can be added if protocol is `HTTPS`",
						Sensitive:   true,
					},
				},
			},
			"default_rule": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"target_group_id": schema.StringAttribute{
						Optional:    true,
						Description: "Client Private Key. Required only if protocol is `HTTPS`",
					},
				},
			},
		},
	}
}

func (l *loadBalancerListenerResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan loadBalancerListenerResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts, err := plan.generateLoadBalancerListenerCreateOpts(ctx)
	if err != nil {
		utils.GeneralError(&response.Diagnostics, ctx, err)
		return
	}

	loadBalancerListener, httpResponse, err := l.Client.CreateLoadBalancerListener(
		ctx,
		plan.LoadBalancerID.ValueString(),
	).LoadBalancerListenerCreateOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state, err := adaptLoadBalancerListenerToLoadBalancerListenerResource(
		*loadBalancerListener,
		ctx,
	)
	if err != nil {
		utils.GeneralError(&response.Diagnostics, ctx, err)
		return
	}

	state.LoadBalancerID = plan.LoadBalancerID
	state.Certificate = plan.Certificate

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (l *loadBalancerListenerResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var state loadBalancerListenerResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	loadBalancerListenerDetails, httpResponse, err := l.Client.GetLoadBalancerListener(
		ctx,
		state.LoadBalancerID.ValueString(),
		state.ListenerID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	newState, err := adaptLoadBalancerListenerDetailsToLoadBalancerListenerResource(*loadBalancerListenerDetails, ctx)
	if err != nil {
		utils.GeneralError(&response.Diagnostics, ctx, err)
		return
	}

	newState.LoadBalancerID = state.LoadBalancerID
	newState.ListenerID = state.ListenerID
	if newState.Certificate.IsNull() {
		newState.Certificate = state.Certificate
	}

	response.Diagnostics.Append(response.State.Set(ctx, newState)...)
}

func (l *loadBalancerListenerResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan loadBalancerListenerResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts, err := plan.generateLoadBalancerListenerUpdateOpts(ctx)
	if err != nil {
		utils.GeneralError(&response.Diagnostics, ctx, err)
		return
	}

	loadBalancerListener, httpResponse, err := l.Client.
		UpdateLoadBalancerListener(
			ctx,
			plan.LoadBalancerID.ValueString(),
			plan.ListenerID.ValueString(),
		).
		LoadBalancerListenerOpts(*opts).
		Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state, err := adaptLoadBalancerListenerToLoadBalancerListenerResource(
		*loadBalancerListener,
		ctx,
	)
	if err != nil {
		utils.GeneralError(&response.Diagnostics, ctx, err)
		return
	}

	if state.Certificate.IsNull() {
		state.Certificate = plan.Certificate
	}
	state.LoadBalancerID = plan.LoadBalancerID

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (l *loadBalancerListenerResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state loadBalancerListenerResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	httpResponse, err := l.Client.DeleteLoadBalancerListener(
		ctx,
		state.LoadBalancerID.ValueString(),
		state.ListenerID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
	}
}

func (l *loadBalancerListenerResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	idParts := strings.Split(request.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		utils.UnexpectedImportIdentifierError(
			&response.Diagnostics,
			"load_balancer_id,listener_id",
			request.ID,
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("load_balancer_id"),
		idParts[0],
	)...)
	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("listener_id"),
		idParts[1],
	)...)
}

func NewLoadBalancerListenerResource() resource.Resource {
	return &loadBalancerListenerResource{
		PubliccloudResourceAPI: utils.PubliccloudResourceAPI{
			Name: "public_cloud_load_balancer_listener",
		},
	}
}
