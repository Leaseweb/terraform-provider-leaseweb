package opts

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
	"testing"
)

func setupSdkInstance(
	instance *publicCloud.Instance,
	operatingSystem *publicCloud.OperatingSystem,
	contract *publicCloud.Contract,
	instanceType *publicCloud.InstanceType,
) *publicCloud.Instance {
	if instance == nil {
		instance = publicCloud.NewInstance()
	}

	if instanceType == nil {
		instanceType, _ = publicCloud.NewInstanceTypeFromValue("lsw.m5a.4xlarge")
	}

	if operatingSystem == nil {
		operatingSystemId, _ := publicCloud.NewOperatingSystemIdFromValue(
			"UBUNTU_24_04_64BIT",
		)
		operatingSystem = publicCloud.NewOperatingSystem()
		operatingSystem.SetId(*operatingSystemId)
	}

	instance.SetOperatingSystem(*operatingSystem)
	instance.SetResources(*publicCloud.NewInstanceResources())
	instance.SetType(*instanceType)

	if contract != nil {
		instance.SetContract(*contract)
	}

	return instance
}

func TestInstanceOpts_setOptionalUpdateInstanceOpts_incorrectInstanceType(t *testing.T) {
	sdkInstance := setupSdkInstance(nil, nil, nil, nil)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())
	instance.Type = basetypes.NewStringValue("tralala")

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	_, err := instanceOpts.NewUpdateInstanceOpts()

	assert.NotNil(t,
		err,
		"setOptionalUpdateInstanceOpts should return an error",
	)
}

func TestInstanceOpts_setOptionalUpdateInstanceOpts(t *testing.T) {
	sdkInstance := publicCloud.NewInstance()
	sdkInstance.SetReference("reference")
	sdkInstance.SetRootDiskSize(23)

	sdkContract := publicCloud.NewContract()
	sdkContract.SetTerm(4)
	sdkContract.SetType("contractType")
	sdkContract.SetBillingFrequency(6)

	sdkInstanceType, _ := publicCloud.NewInstanceTypeFromValue("lsw.m3.xlarge")

	sdkInstance = setupSdkInstance(sdkInstance, nil, sdkContract, sdkInstanceType)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())

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
	sdkInstance := setupSdkInstance(nil, nil, nil, nil)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	_, err := instanceOpts.NewUpdateInstanceOpts()

	assert.Nil(t,
		err,
		"NewUpdateInstanceOpts should not return an error",
	)
}

func TestInstanceOpts_NewUpdateInstanceOpts_error(t *testing.T) {
	sdkInstance := setupSdkInstance(nil, nil, nil, nil)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())
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
	sdkOperatingSystem := publicCloud.NewOperatingSystem()
	sdkOperatingSystem.SetId(*sdkOperatingSystemId)

	sdkContract := publicCloud.NewContract()
	sdkContract.SetTerm(4)
	sdkContract.SetType("contractType")
	sdkContract.SetBillingFrequency(6)

	sdkInstance := publicCloud.NewInstance()
	sdkInstance.SetOperatingSystem(*sdkOperatingSystem)
	sdkInstance.SetRegion("eu-west-1")
	sdkInstance.SetRootDiskStorageType("rootDiskStorageType")

	sdkInstanceType, _ := publicCloud.NewInstanceTypeFromValue("lsw.m3.xlarge")

	sdkInstance = setupSdkInstance(
		sdkInstance,
		sdkOperatingSystem,
		sdkContract,
		sdkInstanceType,
	)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())

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
		"rootDiskStorageType",
		launchInstanceOpts.GetRootDiskStorageType(),
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
	sdkInstance := publicCloud.NewInstance()
	sdkInstance.SetReference("reference")
	sdkInstance.SetRootDiskSize(32)
	sdkInstance.SetMarketAppId("marketAppId")

	sdkInstance = setupSdkInstance(sdkInstance, nil, nil, nil)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())
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
	sdkOperatingSystem := publicCloud.NewOperatingSystem()
	sdkInstance := setupSdkInstance(nil, sdkOperatingSystem, nil, nil)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())

	instanceOpts := NewInstanceOpts(instance, context.TODO())
	_, err := instanceOpts.NewLaunchInstanceOpts()

	assert.NotNil(t, err, "NewLaunchInstanceOpts should return an error")
}

func TestInstanceOpts_NewLaunchInstanceOpts_cannotSetInstanceType(t *testing.T) {
	sdkInstance := setupSdkInstance(nil, nil, nil, nil)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())
	instance.Type = basetypes.NewStringValue("tralala")

	instanceOpts := NewInstanceOpts(instance, context.TODO())
	_, err := instanceOpts.NewLaunchInstanceOpts()

	assert.NotNil(t, err, "NewLaunchInstanceOpts should return an error")
}
