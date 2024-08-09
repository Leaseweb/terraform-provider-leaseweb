package public_cloud_repository

import (
	"context"
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/core/domain/public_cloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/public_cloud_repository/data_adapters/to_domain_entity"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/repositories/public_cloud_repository/data_adapters/to_sdk_model"
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
	publicCLoudAPI       sdk.PublicCloudApi
	token                string
	adaptInstanceDetails func(
		sdkInstance publicCloud.InstanceDetails,
	) (*public_cloud.Instance, error)
	adaptInstance func(
		sdkInstance publicCloud.Instance,
	) (*public_cloud.Instance, error)
	adaptAutoScalingGroupDetails func(
		sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
	) (*public_cloud.AutoScalingGroup, error)
	adaptLoadBalancerDetails func(
		sdkLoadBalancerDetails publicCloud.LoadBalancerDetails,
	) (*public_cloud.LoadBalancer, error)
	adaptToLaunchInstanceOpts func(instance public_cloud.Instance) (
		*publicCloud.LaunchInstanceOpts, error)
	adaptToUpdateInstanceOpts func(instance public_cloud.Instance) (
		*publicCloud.UpdateInstanceOpts, error)
	adaptRegion       func(sdkRegion publicCloud.Region) public_cloud.Region
	adaptInstanceType func(sdkInstanceType publicCloud.InstanceType) (
		*public_cloud.InstanceType,
		error,
	)
	adaptImageDetails func(sdkImage publicCloud.ImageDetails) public_cloud.Image
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
	public_cloud.Instances,
	*shared.RepositoryError,
) {
	var instances public_cloud.Instances

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

		for _, sdkInstance := range result.Instances {
			instance, err := p.adaptInstance(sdkInstance)
			if err != nil {
				return nil, shared.NewGeneralError("GetAllInstances", err)
			}

			instances = append(instances, *instance)
		}

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
) (*public_cloud.Instance, *shared.RepositoryError) {
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

	instance, err := p.adaptInstanceDetails(*sdkInstance)
	if err != nil {
		return nil, shared.NewSdkError(
			fmt.Sprintf("GetInstance %q", id),
			err,
			response,
		)
	}

	return instance, nil
}

func (p PublicCloudRepository) GetAutoScalingGroup(
	id string,
	ctx context.Context,
) (*public_cloud.AutoScalingGroup, *shared.RepositoryError) {
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

	autoScalingGroup, err := p.adaptAutoScalingGroupDetails(
		*sdkAutoScalingGroupDetails,
	)
	if err != nil {
		return nil, shared.NewGeneralError(
			fmt.Sprintf("GetAutoScalingGroup %q", id),
			err,
		)
	}

	return autoScalingGroup, nil
}

func (p PublicCloudRepository) GetLoadBalancer(
	id string,
	ctx context.Context,
) (*public_cloud.LoadBalancer, *shared.RepositoryError) {
	var loadBalancer *public_cloud.LoadBalancer

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

	loadBalancer, err = p.adaptLoadBalancerDetails(*sdkLoadBalancerDetails)
	if err != nil {
		return nil, shared.NewGeneralError(
			fmt.Sprintf("getLoadBalancer %q", sdkLoadBalancerDetails.GetId()),
			err,
		)
	}

	return loadBalancer, nil
}

func (p PublicCloudRepository) CreateInstance(
	instance public_cloud.Instance,
	ctx context.Context,
) (*public_cloud.Instance, *shared.RepositoryError) {

	launchInstanceOpts, err := p.adaptToLaunchInstanceOpts(instance)
	if err != nil {
		return nil, shared.NewGeneralError(
			"CreateInstance",
			err,
		)
	}

	sdkLaunchedInstance, response, err := p.publicCLoudAPI.
		LaunchInstance(p.authContext(ctx)).
		LaunchInstanceOpts(*launchInstanceOpts).Execute()

	if err != nil {
		return nil, shared.NewSdkError(
			"CreateInstance",
			err,
			response,
		)
	}

	launchedInstance, err := p.adaptInstance(*sdkLaunchedInstance)

	if err != nil {
		return nil, shared.NewGeneralError("CreateInstance", err)
	}

	return launchedInstance, nil
}

func (p PublicCloudRepository) UpdateInstance(
	instance public_cloud.Instance,
	ctx context.Context,
) (*public_cloud.Instance, *shared.RepositoryError) {

	updateInstanceOpts, err := p.adaptToUpdateInstanceOpts(instance)
	if err != nil {
		return nil, shared.NewGeneralError(
			fmt.Sprintf("UpdateInstance %q", instance.Id),
			err,
		)
	}

	sdkUpdatedInstance, response, err := p.publicCLoudAPI.UpdateInstance(
		p.authContext(ctx),
		instance.Id,
	).UpdateInstanceOpts(*updateInstanceOpts).Execute()
	if err != nil {
		return nil, shared.NewSdkError(
			fmt.Sprintf("UpdateInstance %q", instance.Id),
			err,
			response,
		)
	}

	updatedInstance, err := p.adaptInstanceDetails(*sdkUpdatedInstance)
	if err != nil {
		return nil, shared.NewGeneralError(
			fmt.Sprintf("UpdateInstance %q", instance.Id),
			err,
		)
	}

	return updatedInstance, nil
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
) (public_cloud.InstanceTypes, *shared.RepositoryError) {
	var instanceTypes public_cloud.InstanceTypes

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

	for _, sdkInstanceType := range sdkInstanceTypes.InstanceTypes {
		instanceType, err := p.adaptInstanceType(sdkInstanceType)
		if err != nil {
			return nil, shared.NewSdkError(
				fmt.Sprintf("GetAvailableInstanceTypesForUpdate %q", id),
				err,
				response,
			)
		}
		instanceTypes = append(instanceTypes, *instanceType)
	}

	return instanceTypes, nil
}

func (p PublicCloudRepository) GetRegions(ctx context.Context) (
	public_cloud.Regions,
	*shared.RepositoryError,
) {
	var regions public_cloud.Regions

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

		for _, sdkRegion := range result.Regions {
			region := p.adaptRegion(sdkRegion)

			regions = append(regions, region)
		}

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
) (public_cloud.InstanceTypes, *shared.RepositoryError) {
	var instanceTypes public_cloud.InstanceTypes

	request := p.publicCLoudAPI.GetInstanceTypeList(p.authContext(ctx)).
		Region(region)

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

		for _, sdkInstanceType := range result.InstanceTypes {
			instanceType, err := p.adaptInstanceType(sdkInstanceType)
			if err != nil {
				return nil, shared.NewSdkError(
					"GetInstanceTypesForRegion",
					err,
					response,
				)
			}

			instanceTypes = append(instanceTypes, *instanceType)
		}

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
	public_cloud.Images,
	*shared.RepositoryError,
) {
	var images public_cloud.Images

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

		for _, sdkImage := range result.Images {
			image := p.adaptImageDetails(sdkImage)
			images = append(images, image)
		}

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
		publicCLoudAPI:               client.PublicCloudAPI,
		token:                        token,
		adaptInstanceDetails:         to_domain_entity.AdaptInstanceDetails,
		adaptInstance:                to_domain_entity.AdaptInstance,
		adaptAutoScalingGroupDetails: to_domain_entity.AdaptAutoScalingGroupDetails,
		adaptLoadBalancerDetails:     to_domain_entity.AdaptLoadBalancerDetails,
		adaptInstanceType:            to_domain_entity.AdaptInstanceType,
		adaptRegion:                  to_domain_entity.AdaptRegion,
		adaptImageDetails:            to_domain_entity.AdaptImageDetails,
		adaptToLaunchInstanceOpts:    to_sdk_model.AdaptToLaunchInstanceOpts,
		adaptToUpdateInstanceOpts:    to_sdk_model.AdaptToUpdateInstanceOpts,
	}
}
