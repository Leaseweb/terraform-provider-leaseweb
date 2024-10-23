package publiccloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

type dataSourceModelContract struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	State            types.String `tfsdk:"state"`
}

func newDataSourceModelContract(sdkContract publicCloud.Contract) dataSourceModelContract {
	return dataSourceModelContract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.GetBillingFrequency())),
		Term:             basetypes.NewInt64Value(int64(sdkContract.GetTerm())),
		Type:             basetypes.NewStringValue(string(sdkContract.GetType())),
		EndsAt:           utils.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.GetState())),
	}
}

type dataSourceModelInstance struct {
	ID                  types.String            `tfsdk:"id"`
	Region              types.String            `tfsdk:"region"`
	Reference           types.String            `tfsdk:"reference"`
	Image               dataSourceModelImage    `tfsdk:"image"`
	State               types.String            `tfsdk:"state"`
	Type                types.String            `tfsdk:"type"`
	RootDiskSize        types.Int64             `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String            `tfsdk:"root_disk_storage_type"`
	Ips                 []dataSourceModelIp     `tfsdk:"ips"`
	Contract            dataSourceModelContract `tfsdk:"contract"`
	MarketAppId         types.String            `tfsdk:"market_app_id"`
}

func adaptSdkInstanceToDatasourceInstance(sdkInstance publicCloud.Instance) dataSourceModelInstance {
	var ips []dataSourceModelIp
	for _, ip := range sdkInstance.Ips {
		ips = append(ips, adaptSdkIpToDatasourceIp(ip))
	}

	return dataSourceModelInstance{
		ID:                  basetypes.NewStringValue(sdkInstance.GetId()),
		Region:              basetypes.NewStringValue(string(sdkInstance.GetRegion())),
		Reference:           basetypes.NewStringPointerValue(sdkInstance.Reference.Get()),
		Image:               adaptSdkImageToDatasourceImage(sdkInstance.GetImage()),
		State:               basetypes.NewStringValue(string(sdkInstance.GetState())),
		Type:                basetypes.NewStringValue(string(sdkInstance.GetType())),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.GetRootDiskSize())),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.GetRootDiskStorageType())),
		Ips:                 ips,
		Contract:            newDataSourceModelContract(sdkInstance.GetContract()),
		MarketAppId:         basetypes.NewStringPointerValue(sdkInstance.MarketAppId.Get()),
	}
}

type dataSourceModelImage struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Custom       types.Bool   `tfsdk:"custom"`
	State        types.String `tfsdk:"state"`
	MarketApps   []string     `tfsdk:"market_apps"`
	StorageTypes []string     `tfsdk:"storage_types"`
	Flavour      types.String `tfsdk:"flavour"`
	Region       types.String `tfsdk:"region"`
}

func adaptSdkImageToDatasourceImage(sdkImage publicCloud.Image) dataSourceModelImage {
	return dataSourceModelImage{
		ID:      basetypes.NewStringValue(sdkImage.GetId()),
		Name:    basetypes.NewStringValue(sdkImage.GetName()),
		Custom:  basetypes.NewBoolValue(sdkImage.GetCustom()),
		Flavour: basetypes.NewStringValue(string(sdkImage.GetFlavour())),
	}
}

func adaptSdkImageDetailsToDatasourceImage(
	sdkImageDetails publicCloud.ImageDetails,
) dataSourceModelImage {
	var marketApps []string
	var storageTypes []string

	for _, marketApp := range sdkImageDetails.GetMarketApps() {
		marketApps = append(marketApps, string(marketApp))
	}

	for _, storageType := range sdkImageDetails.GetStorageTypes() {
		storageTypes = append(storageTypes, string(storageType))
	}

	return dataSourceModelImage{
		ID:           basetypes.NewStringValue(sdkImageDetails.GetId()),
		Name:         basetypes.NewStringValue(sdkImageDetails.GetName()),
		Custom:       basetypes.NewBoolValue(sdkImageDetails.GetCustom()),
		State:        basetypes.NewStringValue(string(sdkImageDetails.GetState())),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue(string(sdkImageDetails.GetFlavour())),
		Region:       basetypes.NewStringValue(string(sdkImageDetails.GetRegion())),
	}
}

type dataSourceModelImages struct {
	Images []dataSourceModelImage `tfsdk:"images"`
}

func adaptSdkImagesToDatasourceImages(sdkImages []publicCloud.ImageDetails) dataSourceModelImages {
	var images dataSourceModelImages

	for _, sdkImageDetails := range sdkImages {
		image := adaptSdkImageDetailsToDatasourceImage(sdkImageDetails)
		images.Images = append(images.Images, image)
	}

	return images
}

type dataSourceModelIp struct {
	Ip types.String `tfsdk:"ip"`
}

func adaptSdkIpToDatasourceIp(sdkIp publicCloud.Ip) dataSourceModelIp {
	return dataSourceModelIp{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}
}

type dataSourceModelInstances struct {
	Instances []dataSourceModelInstance `tfsdk:"instances"`
}

func adaptSdkInstancesToDatasourceInstances(sdkInstances []publicCloud.Instance) dataSourceModelInstances {
	var instances dataSourceModelInstances

	for _, sdkInstance := range sdkInstances {
		instance := adaptSdkInstanceToDatasourceInstance(sdkInstance)
		instances.Instances = append(instances.Instances, instance)
	}

	return instances
}

type reason string

const (
	reasonContractTermCannotBeZero reason = "contract.term cannot be 0 when contract type is MONTHLY"
	reasonContractTermMustBeZero   reason = "contract.term must be 0 when contract type is HOURLY"
	reasonNone                     reason = ""
)

type resourceModelContract struct {
	BillingFrequency types.Int64  `tfsdk:"billing_frequency"`
	Term             types.Int64  `tfsdk:"term"`
	Type             types.String `tfsdk:"type"`
	EndsAt           types.String `tfsdk:"ends_at"`
	State            types.String `tfsdk:"state"`
}

func (c resourceModelContract) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"billing_frequency": types.Int64Type,
		"term":              types.Int64Type,
		"type":              types.StringType,
		"ends_at":           types.StringType,
		"state":             types.StringType,
	}
}

func (c resourceModelContract) IsContractTermValid() (bool, reason) {
	if c.Type.ValueString() == string(publicCloud.CONTRACTTYPE_MONTHLY) && c.Term.ValueInt64() == 0 {
		return false, reasonContractTermCannotBeZero
	}

	if c.Type.ValueString() == string(publicCloud.CONTRACTTYPE_HOURLY) && c.Term.ValueInt64() != 0 {
		return false, reasonContractTermMustBeZero
	}

	return true, reasonNone
}

func adaptSdkContractToResourceContract(
	_ context.Context,
	sdkContract publicCloud.Contract,
) (*resourceModelContract, error) {
	return &resourceModelContract{
		BillingFrequency: basetypes.NewInt64Value(int64(sdkContract.GetBillingFrequency())),
		Term:             basetypes.NewInt64Value(int64(sdkContract.GetTerm())),
		Type:             basetypes.NewStringValue(string(sdkContract.GetType())),
		EndsAt:           utils.AdaptNullableTimeToStringValue(sdkContract.EndsAt.Get()),
		State:            basetypes.NewStringValue(string(sdkContract.GetState())),
	}, nil
}

type reasonInstanceCannotBeTerminated string

type resourceModelInstance struct {
	ID                  types.String `tfsdk:"id"`
	Region              types.String `tfsdk:"region"`
	Reference           types.String `tfsdk:"reference"`
	Image               types.Object `tfsdk:"image"`
	State               types.String `tfsdk:"state"`
	Type                types.String `tfsdk:"type"`
	RootDiskSize        types.Int64  `tfsdk:"root_disk_size"`
	RootDiskStorageType types.String `tfsdk:"root_disk_storage_type"`
	Ips                 types.List   `tfsdk:"ips"`
	Contract            types.Object `tfsdk:"contract"`
	MarketAppId         types.String `tfsdk:"market_app_id"`
}

func (i resourceModelInstance) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":        types.StringType,
		"region":    types.StringType,
		"reference": types.StringType,
		"image": types.ObjectType{
			AttrTypes: resourceModelImage{}.AttributeTypes(),
		},
		"state":                  types.StringType,
		"type":                   types.StringType,
		"root_disk_size":         types.Int64Type,
		"root_disk_storage_type": types.StringType,
		"ips": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: resourceModelIp{}.AttributeTypes(),
			},
		},
		"contract": types.ObjectType{
			AttrTypes: resourceModelContract{}.AttributeTypes(),
		},
		"market_app_id": types.StringType,
	}
}

func (i resourceModelInstance) GetLaunchInstanceOpts(ctx context.Context) (
	*publicCloud.LaunchInstanceOpts,
	error,
) {
	sdkRootDiskStorageType, err := publicCloud.NewStorageTypeFromValue(i.RootDiskStorageType.ValueString())
	if err != nil {
		return nil, err
	}

	image := resourceModelImage{}
	imageDiags := i.Image.As(ctx, &image, basetypes.ObjectAsOptions{})
	if imageDiags != nil {
		return nil, utils.ReturnError("GetLaunchInstanceOpts", imageDiags)
	}

	contract := resourceModelContract{}
	contractDiags := i.Contract.As(ctx, &contract, basetypes.ObjectAsOptions{})
	if contractDiags != nil {
		return nil, utils.ReturnError("GetLaunchInstanceOpts", contractDiags)
	}

	sdkContractType, err := publicCloud.NewContractTypeFromValue(
		contract.Type.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	sdkContractTerm, err := publicCloud.NewContractTermFromValue(
		int32(contract.Term.ValueInt64()),
	)
	if err != nil {
		return nil, err
	}

	sdkBillingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
		int32(contract.BillingFrequency.ValueInt64()),
	)
	if err != nil {
		return nil, err
	}

	sdkRegionName, err := publicCloud.NewRegionNameFromValue(
		i.Region.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	sdkInstanceType, err := publicCloud.NewTypeNameFromValue(
		i.Type.ValueString(),
	)
	if err != nil {
		return nil, err
	}

	opts := publicCloud.NewLaunchInstanceOpts(
		*sdkRegionName,
		*sdkInstanceType,
		image.ID.ValueString(),
		*sdkContractType,
		*sdkContractTerm,
		*sdkBillingFrequency,
		*sdkRootDiskStorageType,
	)

	opts.MarketAppId = i.MarketAppId.ValueStringPointer()
	opts.Reference = i.Reference.ValueStringPointer()
	opts.RootDiskSize = utils.AdaptInt64PointerValueToNullableInt32(i.RootDiskSize)

	return opts, nil
}

func (i resourceModelInstance) GetUpdateInstanceOpts(ctx context.Context) (
	*publicCloud.UpdateInstanceOpts,
	error,
) {
	opts := publicCloud.NewUpdateInstanceOpts()

	opts.Reference = i.Reference.ValueStringPointer()
	opts.RootDiskSize = utils.AdaptInt64PointerValueToNullableInt32(i.RootDiskSize)

	contract := resourceModelContract{}
	diags := i.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if diags.HasError() {
		return nil, utils.ReturnError("GetUpdateInstanceOpts", diags)
	}

	if contract.Type.ValueString() != "" {
		contractType, err := publicCloud.NewContractTypeFromValue(
			contract.Type.ValueString(),
		)
		if err != nil {
			return nil, fmt.Errorf("GetUpdateInstanceOpts: %w", err)
		}
		opts.ContractType = contractType
	}

	if contract.Term.ValueInt64() != 0 {
		contractTerm, err := publicCloud.NewContractTermFromValue(
			int32(contract.Term.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf("GetUpdateInstanceOpts: %w", err)
		}
		opts.ContractTerm = contractTerm
	}

	if contract.BillingFrequency.ValueInt64() != 0 {
		billingFrequency, err := publicCloud.NewBillingFrequencyFromValue(
			int32(contract.BillingFrequency.ValueInt64()),
		)
		if err != nil {
			return nil, fmt.Errorf("GetUpdateInstanceOpts: %w", err)
		}
		opts.BillingFrequency = billingFrequency
	}

	if i.Type.ValueString() != "" {
		instanceType, err := publicCloud.NewTypeNameFromValue(
			i.Type.ValueString(),
		)
		if err != nil {
			return nil, fmt.Errorf("GetUpdateInstanceOpts: %w", err)
		}
		opts.Type = instanceType
	}

	return opts, nil
}

func (i resourceModelInstance) CanBeTerminated(ctx context.Context) *reasonInstanceCannotBeTerminated {
	contract := resourceModelContract{}
	contractDiags := i.Contract.As(
		ctx,
		&contract,
		basetypes.ObjectAsOptions{},
	)
	if contractDiags != nil {
		log.Fatal("cannot convert contract objectType to model")
	}

	if i.State.ValueString() == string(publicCloud.STATE_CREATING) || i.State.ValueString() == string(publicCloud.STATE_DESTROYING) || i.State.ValueString() == string(publicCloud.STATE_DESTROYED) {
		reason := reasonInstanceCannotBeTerminated(
			fmt.Sprintf("state is %q", i.State),
		)

		return &reason
	}

	if !contract.EndsAt.IsNull() {
		reason := reasonInstanceCannotBeTerminated(
			fmt.Sprintf("contract.endsAt is %q", contract.EndsAt.ValueString()),
		)

		return &reason
	}

	return nil
}

func adaptSdkInstanceToResourceInstance(
	sdkInstance publicCloud.Instance,
	ctx context.Context,
) (*resourceModelInstance, error) {
	instance := resourceModelInstance{
		ID:                  basetypes.NewStringValue(sdkInstance.Id),
		Region:              basetypes.NewStringValue(string(sdkInstance.Region)),
		Reference:           basetypes.NewStringPointerValue(sdkInstance.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstance.State)),
		Type:                basetypes.NewStringValue(string(sdkInstance.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstance.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstance.RootDiskStorageType)),
		MarketAppId:         basetypes.NewStringPointerValue(sdkInstance.MarketAppId.Get()),
	}

	image, err := utils.AdaptSdkModelToResourceObject(
		sdkInstance.Image,
		resourceModelImage{}.AttributeTypes(),
		ctx,
		adaptSdkImageToResourceImage,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.Image = image

	ips, err := utils.AdaptSdkModelsToListValue(
		sdkInstance.Ips,
		resourceModelIp{}.AttributeTypes(),
		ctx,
		adaptSdkIpToResourceIp,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.Ips = ips

	contract, err := utils.AdaptSdkModelToResourceObject(
		sdkInstance.Contract,
		resourceModelContract{}.AttributeTypes(),
		ctx,
		adaptSdkContractToResourceContract,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}

func adaptSdkInstanceDetailsToResourceInstance(
	sdkInstanceDetails publicCloud.InstanceDetails,
	ctx context.Context,
) (*resourceModelInstance, error) {
	instance := resourceModelInstance{
		ID:                  basetypes.NewStringValue(sdkInstanceDetails.Id),
		Region:              basetypes.NewStringValue(string(sdkInstanceDetails.Region)),
		Reference:           basetypes.NewStringPointerValue(sdkInstanceDetails.Reference.Get()),
		State:               basetypes.NewStringValue(string(sdkInstanceDetails.State)),
		Type:                basetypes.NewStringValue(string(sdkInstanceDetails.Type)),
		RootDiskSize:        basetypes.NewInt64Value(int64(sdkInstanceDetails.RootDiskSize)),
		RootDiskStorageType: basetypes.NewStringValue(string(sdkInstanceDetails.RootDiskStorageType)),
		MarketAppId:         basetypes.NewStringPointerValue(sdkInstanceDetails.MarketAppId.Get()),
	}

	image, err := utils.AdaptSdkModelToResourceObject(
		sdkInstanceDetails.Image,
		resourceModelImage{}.AttributeTypes(),
		ctx,
		adaptSdkImageToResourceImage,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.Image = image

	ips, err := utils.AdaptSdkModelsToListValue(
		sdkInstanceDetails.Ips,
		resourceModelIp{}.AttributeTypes(),
		ctx,
		adaptSdkIpDetailsToResourceIp,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.Ips = ips

	contract, err := utils.AdaptSdkModelToResourceObject(
		sdkInstanceDetails.Contract,
		resourceModelContract{}.AttributeTypes(),
		ctx,
		adaptSdkContractToResourceContract,
	)
	if err != nil {
		return nil, fmt.Errorf("adaptSdkInstanceToResourceInstance: %w", err)
	}
	instance.Contract = contract

	return &instance, nil
}

type resourceModelIp struct {
	Ip types.String `tfsdk:"ip"`
}

func (i resourceModelIp) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip": types.StringType,
	}
}

func adaptSdkIpToResourceIp(
	_ context.Context,
	sdkIp publicCloud.Ip,
) (*resourceModelIp, error) {
	return &resourceModelIp{
		Ip: basetypes.NewStringValue(sdkIp.Ip),
	}, nil
}

func adaptSdkIpDetailsToResourceIp(
	_ context.Context,
	sdkIpDetails publicCloud.IpDetails,
) (*resourceModelIp, error) {
	return &resourceModelIp{
		Ip: basetypes.NewStringValue(sdkIpDetails.Ip),
	}, nil
}

type resourceModelImage struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Custom       types.Bool   `tfsdk:"custom"`
	State        types.String `tfsdk:"state"`
	MarketApps   types.List   `tfsdk:"market_apps"`
	StorageTypes types.List   `tfsdk:"storage_types"`
	Flavour      types.String `tfsdk:"flavour"`
	Region       types.String `tfsdk:"region"`
}

func (i resourceModelImage) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            types.StringType,
		"name":          types.StringType,
		"custom":        types.BoolType,
		"state":         types.StringType,
		"market_apps":   types.ListType{ElemType: types.StringType},
		"storage_types": types.ListType{ElemType: types.StringType},
		"flavour":       types.StringType,
		"region":        types.StringType,
	}
}

func (i resourceModelImage) GetUpdateImageOpts() publicCloud.UpdateImageOpts {
	return publicCloud.UpdateImageOpts{
		Name: i.Name.ValueString(),
	}
}

func (i resourceModelImage) GetCreateImageOpts() publicCloud.CreateImageOpts {
	return publicCloud.CreateImageOpts{
		Name:       i.Name.ValueString(),
		InstanceId: i.ID.ValueString(),
	}
}

func adaptSdkImageDetailsToResourceImage(
	ctx context.Context,
	sdkImageDetails publicCloud.ImageDetails,
) (*resourceModelImage, error) {
	marketApps, diags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		sdkImageDetails.MarketApps,
	)
	if diags.HasError() {
		return nil, fmt.Errorf(
			diags.Errors()[0].Summary(),
			diags.Errors()[0].Detail(),
		)
	}

	storageTypes, diags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		sdkImageDetails.StorageTypes,
	)
	if diags.HasError() {
		return nil, fmt.Errorf(
			diags.Errors()[0].Summary(),
			diags.Errors()[0].Detail(),
		)
	}

	image := resourceModelImage{
		ID:           basetypes.NewStringValue(sdkImageDetails.GetId()),
		Name:         basetypes.NewStringValue(sdkImageDetails.GetName()),
		Custom:       basetypes.NewBoolValue(sdkImageDetails.GetCustom()),
		State:        basetypes.NewStringValue(string(sdkImageDetails.GetState())),
		MarketApps:   marketApps,
		StorageTypes: storageTypes,
		Flavour:      basetypes.NewStringValue(string(sdkImageDetails.Flavour)),
		Region:       basetypes.NewStringValue(string(sdkImageDetails.GetRegion())),
	}

	return &image, nil
}

func adaptSdkImageToResourceImage(
	_ context.Context,
	sdkImage publicCloud.Image,
) (*resourceModelImage, error) {
	emptyList, _ := basetypes.NewListValue(types.StringType, []attr.Value{})

	return &resourceModelImage{
		ID:           basetypes.NewStringValue(sdkImage.GetId()),
		Name:         basetypes.NewStringValue(sdkImage.GetName()),
		Custom:       basetypes.NewBoolValue(sdkImage.GetCustom()),
		Flavour:      basetypes.NewStringValue(string(sdkImage.GetFlavour())),
		MarketApps:   emptyList,
		StorageTypes: emptyList,
	}, nil
}
