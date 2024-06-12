package opts

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
	"testing"
)

func setupSdkInstance(
	instance *publicCloud.Instance,
	operatingSystem *publicCloud.OperatingSystem,
	contract *publicCloud.Contract,
) *publicCloud.Instance {
	if instance == nil {
		instance = publicCloud.NewInstance()
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

	if contract != nil {
		instance.SetContract(*contract)
	}

	return instance
}

func TestInstanceOpts_NewLaunchInstanceOpts_RequiredValues(t *testing.T) {
	sdkOperatingSystemId, _ := publicCloud.NewOperatingSystemIdFromValue(
		"UBUNTU_24_04_64BIT",
	)
	sdkOperatingSystem := publicCloud.NewOperatingSystem()
	sdkOperatingSystem.SetId(*sdkOperatingSystemId)

	sdkContract := publicCloud.NewContract()
	sdkContract.SetTerm(4)
	sdkContract.SetType("contractType")

	sdkInstance := publicCloud.NewInstance()
	sdkInstance.SetOperatingSystem(*publicCloud.NewOperatingSystem())
	sdkInstance.SetOperatingSystem(*sdkOperatingSystem)
	sdkInstance.SetRegion("eu-west-1")
	sdkInstance.SetRootDiskStorageType("rootDiskStorage")
	sdkInstance.SetType("type")

	sdkInstance = setupSdkInstance(sdkInstance, sdkOperatingSystem, sdkContract)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	launchInstanceOpts, err := instanceOpts.NewLaunchInstanceOpts(&diag.Diagnostics{})

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
		"UBUNTU_24_04_64BIT",
		string(launchInstanceOpts.GetOperatingSystemId()),
		"operating_system id  should be UBUNTU_24_04_64BIT",
	)
	assert.Equal(
		t,
		"contractType",
		launchInstanceOpts.GetContractType(),
		"contract type should be contractType",
	)
	assert.Equal(
		t,
		int32(4),
		launchInstanceOpts.GetContractTerm(),
		"contract term should be 4",
	)
	assert.Equal(
		t,
		"rootDiskStorage",
		launchInstanceOpts.GetRootDiskStorageType(),
		"rootDiskStorageType should be rootDiskStorage",
	)
	assert.Equal(
		t,
		"type",
		launchInstanceOpts.GetType(),
		"type should be type",
	)

	assert.Equal(
		t,
		"",
		launchInstanceOpts.GetMarketAppId(),
		"marketAppId should be empty",
	)
	assert.Equal(
		t,
		"",
		launchInstanceOpts.GetReference(),
		"reference should be empty",
	)
	assert.Equal(
		t,
		int32(0),
		launchInstanceOpts.GetRootDiskSize(),
		"rootDiskSize should be empty",
	)
	assert.Equal(
		t,
		"",
		launchInstanceOpts.GetSshKey(),
		"sshKey should be empty",
	)
}

func TestInstanceOpts_NewLaunchInstanceOpts_OptionalValues(t *testing.T) {
	sdkInstance := publicCloud.NewInstance()
	sdkInstance.SetMarketAppId("marketAppId")
	sdkInstance.SetReference("reference")
	sdkInstance.SetRootDiskSize(23)

	sdkInstance = setupSdkInstance(sdkInstance, nil, nil)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())
	instance.SshKey = types.StringValue("sshKey")

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	launchInstanceOpts, err := instanceOpts.NewLaunchInstanceOpts(&diag.Diagnostics{})

	assert.Nil(t,
		err,
		"NewLaunchInstanceOpts should not return an error",
	)
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
		int32(23),
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

func TestInstanceOpts_NewUpdateInstanceOpts_RequiredValues(t *testing.T) {
	sdkInstance := setupSdkInstance(nil, nil, nil)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	updateInstanceOpts := instanceOpts.NewUpdateInstanceOpts()

	assert.Equal(
		t,
		"",
		updateInstanceOpts.GetType(),
		"type should be empty",
	)
	assert.Equal(
		t,
		"",
		updateInstanceOpts.GetReference(),
		"reference should be empty",
	)
	assert.Equal(
		t,
		int32(0),
		updateInstanceOpts.GetRootDiskSize(),
		"rootDiskSize should be empty",
	)
	assert.Equal(
		t,
		"",
		updateInstanceOpts.GetContractType(),
		"contractType should be empty",
	)
	assert.Equal(
		t,
		int32(0),
		updateInstanceOpts.GetContractTerm(),
		"contractTerm should be empty",
	)
	assert.Equal(
		t,
		int32(0),
		updateInstanceOpts.GetBillingFrequency(),
		"billingFrequency should be empty",
	)
}

func TestInstanceOpts_NewUpdateInstanceOpts_OptionalValues(t *testing.T) {
	sdkInstance := publicCloud.NewInstance()
	sdkInstance.SetType("type")
	sdkInstance.SetReference("reference")
	sdkInstance.SetRootDiskSize(32)

	sdkContract := publicCloud.NewContract()
	sdkContract.SetType("contractType")
	sdkContract.SetTerm(4)
	sdkContract.SetBillingFrequency(5)

	sdkInstance = setupSdkInstance(sdkInstance, nil, sdkContract)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	updateInstanceOpts := instanceOpts.NewUpdateInstanceOpts()

	assert.Equal(
		t,
		"type",
		updateInstanceOpts.GetType(),
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
		int32(32),
		updateInstanceOpts.GetRootDiskSize(),
		"rootDiskSize should be set",
	)
	assert.Equal(
		t,
		"contractType",
		updateInstanceOpts.GetContractType(),
		"contractType should be set",
	)
	assert.Equal(
		t,
		int32(4),
		updateInstanceOpts.GetContractTerm(),
		"contractTerm should be set",
	)
	assert.Equal(
		t,
		int32(5),
		updateInstanceOpts.GetBillingFrequency(),
		"billingFrequency should be set",
	)
}
