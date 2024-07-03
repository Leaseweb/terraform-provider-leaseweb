package opts

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
)

func setupSdkInstanceDetails(
	instanceDetails *publicCloud.InstanceDetails,
	operatingSystemDetails *publicCloud.OperatingSystemDetails,
	contract *publicCloud.Contract,
	instanceTypeName *publicCloud.InstanceTypeName,
) *publicCloud.InstanceDetails {
	if instanceDetails == nil {
		instanceDetails = &publicCloud.InstanceDetails{}
	}

	if instanceTypeName == nil {
		instanceTypeName, _ = publicCloud.NewInstanceTypeNameFromValue("lsw.m5a.4xlarge")
	}

	if operatingSystemDetails == nil {
		operatingSystemId, _ := publicCloud.NewOperatingSystemIdFromValue(
			"UBUNTU_24_04_64BIT",
		)
		operatingSystemDetails = &publicCloud.OperatingSystemDetails{}
		operatingSystemDetails.SetId(*operatingSystemId)
	}

	instanceDetails.SetOperatingSystem(*operatingSystemDetails)
	instanceDetails.SetResources(publicCloud.Resources{})
	instanceDetails.SetType(*instanceTypeName)

	if contract != nil {
		instanceDetails.SetContract(*contract)
	}

	return instanceDetails
}

func TestInstanceOpts_setOptionalUpdateInstanceOpts_incorrectInstanceType(t *testing.T) {
	sdkInstanceDetails := setupSdkInstanceDetails(
		nil,
		nil,
		nil,
		nil,
	)

	instance := model.Instance{}
	instance.Populate(sdkInstanceDetails, nil, context.TODO())
	instance.Type = basetypes.NewStringValue("tralala")

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	_, err := instanceOpts.NewUpdateInstanceOpts()

	assert.NotNil(t,
		err,
		"setOptionalUpdateInstanceOpts should return an error",
	)
}

func TestInstanceOpts_setOptionalUpdateInstanceOpts(t *testing.T) {
	sdkInstanceDetails := publicCloud.InstanceDetails{}
	sdkInstanceDetails.SetReference("reference")
	sdkInstanceDetails.SetRootDiskSize(23)

	sdkContract := publicCloud.Contract{}
	sdkContract.SetTerm(4)
	sdkContract.SetType("contractType")
	sdkContract.SetBillingFrequency(6)

	sdkInstanceTypeName, _ := publicCloud.NewInstanceTypeNameFromValue("lsw.m3.xlarge")

	sdkInstanceDetails = *setupSdkInstanceDetails(
		&sdkInstanceDetails,
		nil,
		&sdkContract,
		sdkInstanceTypeName,
	)

	instance := model.Instance{}
	instance.Populate(&sdkInstanceDetails, nil, context.TODO())

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	updateInstanceOpts, err := instanceOpts.NewUpdateInstanceOpts()

	assert.Nil(t,
		err,
		"setOptionalUpdateInstanceOpts should not return an error",
	)
	assert.Equal(
		t,
		"lsw.m3.xlarge",
		string(updateInstanceOpts.GetType()),
		"type should be set",
	)
	assert.Equal(
		t,
		"reference",
		updateInstanceOpts.GetReference(),
		"reference should be set",
	)
	assert.Equal(
		t,
		int32(23),
		updateInstanceOpts.GetRootDiskSize(),
		"rootDiskSize should be set",
	)

	assert.Equal(
		t,
		"contractType",
		updateInstanceOpts.GetContractType(),
		"contract.type should be contractType",
	)
	assert.Equal(
		t,
		int32(4),
		updateInstanceOpts.GetContractTerm(),
		"contract.term should be 4",
	)
	assert.Equal(
		t,
		int32(6),
		updateInstanceOpts.GetBillingFrequency(),
		"contract.billing_frequency should be 6",
	)
}

func TestInstanceOpts_NewUpdateInstanceOpts(t *testing.T) {
	sdkInstance := setupSdkInstanceDetails(
		nil,
		nil,
		nil,
		nil,
	)

	instance := model.Instance{}
	instance.Populate(sdkInstance, nil, context.TODO())

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	_, err := instanceOpts.NewUpdateInstanceOpts()

	assert.Nil(t,
		err,
		"NewUpdateInstanceOpts should not return an error",
	)
}

func TestInstanceOpts_NewUpdateInstanceOpts_error(t *testing.T) {
	sdkInstance := setupSdkInstanceDetails(
		nil,
		nil,
		nil,
		nil,
	)

	instance := model.Instance{}
	instance.Populate(sdkInstance, nil, context.TODO())
	instance.Type = basetypes.NewStringValue("tralala")

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	_, err := instanceOpts.NewUpdateInstanceOpts()

	assert.NotNil(t,
		err,
		"NewUpdateInstanceOpts should return an error",
	)

}

func TestInstanceOpts_NewLaunchInstanceOpts(t *testing.T) {
	sdkOperatingSystemId, _ := publicCloud.NewOperatingSystemIdFromValue(
		"UBUNTU_24_04_64BIT",
	)
	sdkOperatingSystemDetails := publicCloud.OperatingSystemDetails{Id: *sdkOperatingSystemId}
	sdkContract := publicCloud.Contract{
		Term:             4,
		Type:             "contractType",
		BillingFrequency: 6,
	}
	rootDiskStorageType, _ := publicCloud.NewRootDiskStorageTypeFromValue("CENTRAL")

	sdkInstanceDetails := publicCloud.InstanceDetails{
		OperatingSystem:     sdkOperatingSystemDetails,
		Region:              "eu-west-1",
		RootDiskStorageType: *rootDiskStorageType,
	}

	sdkInstanceTypeName, _ := publicCloud.NewInstanceTypeNameFromValue("lsw.m3.xlarge")

	sdkInstanceDetails = *setupSdkInstanceDetails(
		&sdkInstanceDetails,
		&sdkOperatingSystemDetails,
		&sdkContract,
		sdkInstanceTypeName,
	)

	instance := model.Instance{}
	instance.Populate(&sdkInstanceDetails, nil, context.TODO())

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	launchInstanceOpts, err := instanceOpts.NewLaunchInstanceOpts()

	assert.Nil(t,
		err,
		"NewLaunchInstanceOpts should not return an error",
	)
	assert.Equal(
		t,
		"eu-west-1",
		launchInstanceOpts.GetRegion(),
		"region should be eu-west-1",
	)
	assert.Equal(
		t,
		"CENTRAL",
		string(launchInstanceOpts.GetRootDiskStorageType()),
		"rootDiskStorageType should be rootDiskStorageType",
	)
	assert.Equal(
		t,
		"lsw.m3.xlarge",
		string(launchInstanceOpts.GetType()),
		"type should be lsw.m3.xlarge",
	)

	assert.Equal(
		t,
		"UBUNTU_24_04_64BIT",
		string(launchInstanceOpts.GetOperatingSystemId()),
		"operating_system id  should be UBUNTU_24_04_64BIT",
	)

	assert.Equal(
		t,
		"contractType",
		launchInstanceOpts.GetContractType(),
		"contract.type should be contractType",
	)
	assert.Equal(
		t,
		int32(4),
		launchInstanceOpts.GetContractTerm(),
		"contract.term should be 4",
	)
	assert.Equal(
		t,
		int32(6),
		launchInstanceOpts.GetBillingFrequency(),
		"contract.billing_frequency should be 6",
	)
}

func TestInstanceOpts_setOptionalLaunchInstanceOpts(t *testing.T) {
	sdkInstanceDetails := publicCloud.InstanceDetails{}
	sdkInstanceDetails.SetReference("reference")
	sdkInstanceDetails.SetRootDiskSize(32)
	sdkInstanceDetails.SetMarketAppId("marketAppId")

	sdkInstanceDetails = *setupSdkInstanceDetails(
		&sdkInstanceDetails,
		nil,
		nil,
		nil,
	)

	instance := model.Instance{}
	instance.Populate(&sdkInstanceDetails, nil, context.TODO())
	instance.SshKey = types.StringValue("sshKey")

	launchInstanceOpts := publicCloud.LaunchInstanceOpts{}

	instanceOpts := NewInstanceOpts(instance, context.TODO())
	instanceOpts.setOptionalLaunchInstanceOpts(&launchInstanceOpts)

	assert.Equal(
		t,
		"marketAppId",
		launchInstanceOpts.GetMarketAppId(),
		"marketAppId should be set",
	)
	assert.Equal(
		t,
		"reference",
		launchInstanceOpts.GetReference(),
		"reference should be set",
	)
	assert.Equal(
		t,
		int32(32),
		launchInstanceOpts.GetRootDiskSize(),
		"rootDiskSize should be set",
	)
	assert.Equal(
		t,
		"sshKey",
		launchInstanceOpts.GetSshKey(),
		"sshKey should be set",
	)
}

func TestInstanceOpts_NewLaunchInstanceOpts_cannotSetOperatingSystemId(t *testing.T) {
	sdkOperatingSystemDetails := publicCloud.OperatingSystemDetails{}
	sdkInstance := setupSdkInstanceDetails(
		nil,
		&sdkOperatingSystemDetails,
		nil,
		nil,
	)

	instance := model.Instance{}
	instance.Populate(sdkInstance, nil, context.TODO())

	instanceOpts := NewInstanceOpts(instance, context.TODO())
	_, err := instanceOpts.NewLaunchInstanceOpts()

	assert.NotNil(t, err, "NewLaunchInstanceOpts should return an error")
}

func TestInstanceOpts_NewLaunchInstanceOpts_cannotSetInstanceType(t *testing.T) {
	sdkInstance := setupSdkInstanceDetails(
		nil,
		nil,
		nil,
		nil,
	)

	instance := model.Instance{}
	instance.Populate(sdkInstance, nil, context.TODO())
	instance.Type = basetypes.NewStringValue("tralala")

	instanceOpts := NewInstanceOpts(instance, context.TODO())
	_, err := instanceOpts.NewLaunchInstanceOpts()

	assert.NotNil(t, err, "NewLaunchInstanceOpts should return an error")
}

func TestInstanceOpts_NewLaunchInstanceOpts_cannotSetRootDiskStorageType(t *testing.T) {
	sdkInstance := setupSdkInstanceDetails(
		nil,
		nil,
		nil,
		nil,
	)

	instance := model.Instance{}
	instance.Populate(sdkInstance, nil, context.TODO())
	instance.RootDiskStorageType = basetypes.NewStringValue("tralala")

	instanceOpts := NewInstanceOpts(instance, context.TODO())
	_, err := instanceOpts.NewLaunchInstanceOpts()

	assert.NotNil(t, err, "NewLaunchInstanceOpts should return an error")
}
