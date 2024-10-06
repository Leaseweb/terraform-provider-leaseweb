// Package public_cloud_repository implements repository logic
// to access the public_cloud sdk.
package public_cloud_repository

import (
	"context"
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/sdk"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/shared"
)

// Optional contains optional values that can be passed to NewPublicCloudRepository.
type Optional struct {
	Host   *string
	Scheme *string
}

// PublicCloudRepository fulfills contract for ports.PublicCloudRepository.
type PublicCloudRepository struct {
	publicCLoudAPI sdk.PublicCloudApi
	token          string
}

// Injects the authentication token into the context for the sdk.
func (p PublicCloudRepository) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		publicCloud.ContextAPIKeys,
		map[string]publicCloud.APIKey{
			"X-LSW-Auth": {Key: p.token, Prefix: ""},
		},
	)
}

func (p PublicCloudRepository) GetAllInstances(ctx context.Context) (
	[]publicCloud.GetInstanceListResult,
	*shared.RepositoryError,
) {
	var instances []publicCloud.GetInstanceListResult

	request := p.publicCLoudAPI.GetInstanceList(p.authContext(ctx))

	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError("GetAllInstances", err, response)
	}

	metadata := result.GetMetadata()
	pagination := shared.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, shared.NewSdkError("GetAllInstances", err, response)
		}

		instances = append(instances, *result)

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, shared.NewSdkError("GetAllInstances", err, response)
		}
	}

	return instances, nil
}

func (p PublicCloudRepository) GetInstance(
	id string,
	ctx context.Context,
) (*publicCloud.InstanceDetails, *shared.RepositoryError) {
	sdkInstance, response, err := p.publicCLoudAPI.GetInstance(
		p.authContext(ctx),
		id,
	).Execute()

	if err != nil {
		return nil, shared.NewSdkError(
			fmt.Sprintf("GetInstance %q", id),
			err,
			response,
		)
	}

	return sdkInstance, nil
}

func (p PublicCloudRepository) GetAutoScalingGroup(
	id string,
	ctx context.Context,
) (*publicCloud.AutoScalingGroupDetails, *shared.RepositoryError) {
	sdkAutoScalingGroupDetails, response, err := p.publicCLoudAPI.GetAutoScalingGroup(
		p.authContext(ctx),
		id,
	).Execute()
	if err != nil {
		return nil, shared.NewSdkError(
			fmt.Sprintf("GetAutoScalingGroup %q", id),
			err,
			response,
		)
	}

	return sdkAutoScalingGroupDetails, nil
}

func (p PublicCloudRepository) GetLoadBalancer(
	id string,
	ctx context.Context,
) (*publicCloud.LoadBalancerDetails, *shared.RepositoryError) {
	sdkLoadBalancerDetails, response, err := p.publicCLoudAPI.GetLoadBalancer(
		p.authContext(ctx),
		id,
	).Execute()
	if err != nil {
		return nil, shared.NewSdkError(
			fmt.Sprintf("getLoadBalancer %q", id),
			err,
			response,
		)
	}

	return sdkLoadBalancerDetails, nil
}

func (p PublicCloudRepository) CreateInstance(
	opts publicCloud.LaunchInstanceOpts,
	ctx context.Context,
) (*publicCloud.Instance, *shared.RepositoryError) {

	sdkLaunchedInstance, response, err := p.publicCLoudAPI.
		LaunchInstance(p.authContext(ctx)).
		LaunchInstanceOpts(opts).Execute()

	if err != nil {
		return nil, shared.NewSdkError(
			"CreateInstance",
			err,
			response,
		)
	}

	return sdkLaunchedInstance, nil
}

func (p PublicCloudRepository) UpdateInstance(
	opts publicCloud.UpdateInstanceOpts,
	id string,
	ctx context.Context,
) (*publicCloud.InstanceDetails, *shared.RepositoryError) {

	sdkUpdatedInstance, response, err := p.publicCLoudAPI.UpdateInstance(
		p.authContext(ctx),
		id,
	).UpdateInstanceOpts(opts).Execute()
	if err != nil {
		return nil, shared.NewSdkError(
			fmt.Sprintf("UpdateInstance %q", id),
			err,
			response,
		)
	}

	return sdkUpdatedInstance, nil
}

func (p PublicCloudRepository) DeleteInstance(
	id string,
	ctx context.Context,
) *shared.RepositoryError {
	response, err := p.publicCLoudAPI.TerminateInstance(
		p.authContext(ctx),
		id,
	).Execute()
	if err != nil {
		return shared.NewSdkError(
			fmt.Sprintf("DeleteInstance %q", id),
			err,
			response,
		)
	}

	return nil
}

func (p PublicCloudRepository) GetAvailableInstanceTypesForUpdate(
	id string,
	ctx context.Context,
) ([]publicCloud.InstanceTypes, *shared.RepositoryError) {
	var instanceTypes []publicCloud.InstanceTypes

	sdkInstanceTypes, response, err := p.publicCLoudAPI.GetUpdateInstanceTypeList(
		p.authContext(ctx),
		id,
	).Execute()
	if err != nil {
		return nil, shared.NewSdkError(
			fmt.Sprintf("GetAvailableInstanceTypesForUpdate %q", id),
			err,
			response,
		)
	}

	instanceTypes = append(instanceTypes, *sdkInstanceTypes)
	return instanceTypes, nil
}

func (p PublicCloudRepository) GetRegions(ctx context.Context) (
	[]publicCloud.GetRegionListResult,
	*shared.RepositoryError,
) {
	var regions []publicCloud.GetRegionListResult

	request := p.publicCLoudAPI.GetRegionList(p.authContext(ctx))

	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError("GetRegions", err, response)
	}

	metadata := result.GetMetadata()
	pagination := shared.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, shared.NewSdkError("GetRegions", err, response)
		}

		regions = append(regions, *result)

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, shared.NewSdkError("GetAllInstances", err, response)
		}
	}

	return regions, nil
}

func (p PublicCloudRepository) GetInstanceTypesForRegion(
	region string,
	ctx context.Context,
) ([]publicCloud.InstanceTypes, *shared.RepositoryError) {
	var instanceTypes []publicCloud.InstanceTypes

	request := p.publicCLoudAPI.GetInstanceTypeList(p.authContext(ctx)).
		Region(publicCloud.RegionName(region))

	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError(
			"GetInstanceTypesForRegion",
			err,
			response,
		)
	}

	metadata := result.GetMetadata()
	pagination := shared.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, shared.NewSdkError(
				"GetInstanceTypesForRegion",
				err,
				response,
			)
		}

		instanceTypes = append(instanceTypes, *result)

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, shared.NewSdkError("GetAllInstances", err, response)
		}
	}

	return instanceTypes, nil
}

func (p PublicCloudRepository) GetAllImages(ctx context.Context) (
	[]publicCloud.GetImageListResult,
	*shared.RepositoryError,
) {
	var images []publicCloud.GetImageListResult

	request := p.publicCLoudAPI.GetImageList(p.authContext(ctx))

	result, response, err := request.Execute()

	if err != nil {
		return nil, shared.NewSdkError("GetAllImages", err, response)
	}

	metadata := result.GetMetadata()
	pagination := shared.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, shared.NewSdkError("GetAllImages", err, response)
		}

		images = append(images, *result)

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, shared.NewSdkError("GetAllImages", err, response)
		}
	}

	return images, nil
}

func NewPublicCloudRepository(
	token string,
	optional Optional,
) PublicCloudRepository {
	configuration := publicCloud.NewConfiguration()

	if optional.Host != nil {
		configuration.Host = *optional.Host
	}
	if optional.Scheme != nil {
		configuration.Scheme = *optional.Scheme
	}

	client := *publicCloud.NewAPIClient(configuration)

	return PublicCloudRepository{
		publicCLoudAPI: client.PublicCloudAPI,
		token:          token,
	}
}
