package public_cloud_repository

import (
	"context"
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/repositories/public_cloud_repository/data_adapters/to_instance"
	"terraform-provider-leaseweb/internal/repositories/public_cloud_repository/data_adapters/to_sdk_model"
	"terraform-provider-leaseweb/internal/repositories/sdk"
	"terraform-provider-leaseweb/internal/repositories/shared"
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
	) (*domain.Instance, error)
	adaptInstance func(
		sdkInstance publicCloud.Instance,
	) (*domain.Instance, error)
	adaptAutoScalingGroupDetails func(
		sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
	) (*domain.AutoScalingGroup, error)
	adaptLoadBalancerDetails func(
		sdkLoadBalancerDetails publicCloud.LoadBalancerDetails,
	) (*domain.LoadBalancer, error)
	adaptToLaunchInstanceOpts func(instance domain.Instance) (
		*publicCloud.LaunchInstanceOpts, error)
	adaptToUpdateInstanceOpts func(instance domain.Instance) (
		*publicCloud.UpdateInstanceOpts, error)
	adaptRegion       func(sdkRegion publicCloud.Region) domain.Region
	adaptInstanceType func(sdkInstanceType publicCloud.InstanceType) (
		*domain.InstanceType,
		error,
	)
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
	domain.Instances,
	*shared.RepositoryError,
) {
	var instances domain.Instances

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
	id value_object.Uuid,
	ctx context.Context,
) (*domain.Instance, *shared.RepositoryError) {
	sdkInstance, response, err := p.publicCLoudAPI.GetInstance(
		p.authContext(ctx),
		id.String(),
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
	id value_object.Uuid,
	ctx context.Context,
) (*domain.AutoScalingGroup, *shared.RepositoryError) {
	sdkAutoScalingGroupDetails, response, err := p.publicCLoudAPI.GetAutoScalingGroup(
		p.authContext(ctx),
		id.String(),
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
	id value_object.Uuid,
	ctx context.Context,
) (*domain.LoadBalancer, *shared.RepositoryError) {
	var loadBalancer *domain.LoadBalancer

	sdkLoadBalancerDetails, response, err := p.publicCLoudAPI.GetLoadBalancer(
		p.authContext(ctx),
		id.String(),
	).Execute()
	if err != nil {
		return nil, shared.NewSdkError(
			fmt.Sprintf("GetLoadBalancer %q", id),
			err,
			response,
		)
	}

	loadBalancer, err = p.adaptLoadBalancerDetails(*sdkLoadBalancerDetails)
	if err != nil {
		return nil, shared.NewGeneralError(
			fmt.Sprintf("GetLoadBalancer %q", sdkLoadBalancerDetails.GetId()),
			err,
		)
	}

	return loadBalancer, nil
}

func (p PublicCloudRepository) CreateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, *shared.RepositoryError) {

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
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, *shared.RepositoryError) {

	updateInstanceOpts, err := p.adaptToUpdateInstanceOpts(instance)
	if err != nil {
		return nil, shared.NewGeneralError(
			fmt.Sprintf("UpdateInstance %q", instance.Id),
			err,
		)
	}

	sdkUpdatedInstance, response, err := p.publicCLoudAPI.UpdateInstance(
		p.authContext(ctx),
		instance.Id.String(),
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
	id value_object.Uuid,
	ctx context.Context,
) *shared.RepositoryError {
	response, err := p.publicCLoudAPI.TerminateInstance(
		p.authContext(ctx),
		id.String(),
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
	id value_object.Uuid,
	ctx context.Context,
) (domain.InstanceTypes, *shared.RepositoryError) {
	var instanceTypes domain.InstanceTypes

	sdkInstanceTypes, response, err := p.publicCLoudAPI.GetUpdateInstanceTypeList(
		p.authContext(ctx),
		id.String(),
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
	domain.Regions,
	*shared.RepositoryError,
) {
	var regions domain.Regions

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
) (domain.InstanceTypes, *shared.RepositoryError) {
	var instanceTypes domain.InstanceTypes

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
		adaptInstanceDetails:         to_sdk_model.AdaptInstanceDetails,
		adaptInstance:                to_sdk_model.AdaptInstance,
		adaptAutoScalingGroupDetails: to_sdk_model.AdaptAutoScalingGroupDetails,
		adaptLoadBalancerDetails:     to_sdk_model.AdaptLoadBalancerDetails,
		adaptInstanceType:            to_sdk_model.AdaptInstanceType,
		adaptRegion:                  to_sdk_model.AdaptRegion,
		adaptToLaunchInstanceOpts:    to_instance.AdaptToLaunchInstanceOpts,
		adaptToUpdateInstanceOpts:    to_instance.AdaptToUpdateInstanceOpts,
	}
}
