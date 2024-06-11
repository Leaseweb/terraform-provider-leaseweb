package opts

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/public_cloud/resource/instance/model"
	"testing"
)

func TestInstanceOpts_NewLaunchInstanceOpts_RequiredValues(t *testing.T) {
	sdkContract := publicCloud.NewContract()
	sdkContract.SetTerm(4)
	sdkContract.SetType("contractType")

	sdkOperatingSystem := publicCloud.NewOperatingSystem()
	sdkOperatingSystem.SetId("operatingSystemId")

	sdkInstance := publicCloud.NewInstance()
	sdkInstance.SetOperatingSystem(*publicCloud.NewOperatingSystem())
	sdkInstance.SetResources(*publicCloud.NewInstanceResources())
	sdkInstance.SetContract(*sdkContract)
	sdkInstance.SetOperatingSystem(*sdkOperatingSystem)
	sdkInstance.SetRegion("eu-west-1")
	sdkInstance.SetRootDiskStorageType("rootDiskStorage")

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	launchInstanceOpts := instanceOpts.NewLaunchInstanceOpts()

	assert.Equal(
		t,
		"eu-west-1",
		launchInstanceOpts.GetRegion(),
		"region should be eu-west-1",
	)
	assert.Equal(
		t,
		"\"operatingSystemId\"",
		launchInstanceOpts.GetOperatingSystemId(),
		"operating_system id  should be operatingSystemId",
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
		"",
		launchInstanceOpts.GetType(),
		"type should be empty",
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
	sdkInstance.SetOperatingSystem(*publicCloud.NewOperatingSystem())
	sdkInstance.SetResources(*publicCloud.NewInstanceResources())

	sdkInstance.SetType("type")
	sdkInstance.SetMarketAppId("marketAppId")
	sdkInstance.SetReference("reference")
	sdkInstance.SetRootDiskSize(23)

	instance := model.Instance{}
	instance.Populate(sdkInstance, context.TODO())
	instance.SshKey = types.StringValue("sshKey")

	instanceOpts := NewInstanceOpts(instance, context.TODO())

	launchInstanceOpts := instanceOpts.NewLaunchInstanceOpts()

	assert.Equal(
		t,
		"type",
		launchInstanceOpts.GetType(),
		"type should be set",
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
	sdkInstance := publicCloud.NewInstance()
	sdkInstance.SetOperatingSystem(*publicCloud.NewOperatingSystem())
	sdkInstance.SetResources(*publicCloud.NewInstanceResources())

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
	sdkContract := publicCloud.NewContract()
	sdkContract.SetType("contractType")
	sdkContract.SetTerm(4)
	sdkContract.SetBillingFrequency(5)

	sdkInstance := publicCloud.NewInstance()
	sdkInstance.SetType("type")
	sdkInstance.SetReference("reference")
	sdkInstance.SetRootDiskSize(32)
	sdkInstance.SetOperatingSystem(*publicCloud.NewOperatingSystem())
	sdkInstance.SetResources(*publicCloud.NewInstanceResources())
	sdkInstance.SetContract(*sdkContract)

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
