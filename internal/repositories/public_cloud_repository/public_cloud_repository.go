package public_cloud_repository

import (
	"context"
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
	"terraform-provider-leaseweb/internal/repositories/shared"
)

type Optional struct {
	Host   *string
	Scheme *string
}

type PublicCloudRepository struct {
	publicCLoudAPI         publicCloudApi
	token                  string
	convertInstanceDetails func(
		sdkInstance publicCloud.InstanceDetails,
	) (*domain.Instance, error)
	convertInstance func(sdkInstance publicCloud.Instance) (
		*domain.Instance,
		error,
	)
	convertAutoScalingGroupDetails func(
		sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
	) (*domain.AutoScalingGroup, error)
	convertLoadBalancerDetails func(
		sdkLoadBalancerDetails publicCloud.LoadBalancerDetails,
	) (*domain.LoadBalancer, error)
	convertEntityToLaunchInstanceOpts func(instance domain.Instance) (
		*publicCloud.LaunchInstanceOpts, error)
	convertEntityToUpdateInstanceOpts func(instance domain.Instance) (
		*publicCloud.UpdateInstanceOpts, error)
	convertRegion       func(sdkRegion publicCloud.Region) domain.Region
	convertInstanceType func(sdkInstanceType publicCloud.InstanceType) domain.InstanceType
}

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

	result, response, err := p.publicCLoudAPI.GetInstanceList(p.authContext(ctx)).Execute()

	if err != nil {
		return nil, shared.NewSdkError("GetAllInstances", err, response)
	}

	for _, sdkInstance := range result.Instances {
		instance, err := p.convertInstance(sdkInstance)
		if err != nil {
			return nil, shared.NewGeneralError("GetAllInstances", err)
		}

		instances = append(instances, *instance)
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

	instance, err := p.convertInstanceDetails(*sdkInstance)
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

	autoScalingGroup, err := p.convertAutoScalingGroupDetails(*sdkAutoScalingGroupDetails)
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

	loadBalancer, err = p.convertLoadBalancerDetails(*sdkLoadBalancerDetails)
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

	launchInstanceOpts, err := p.convertEntityToLaunchInstanceOpts(instance)
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

	launchedInstance, err := p.convertInstance(*sdkLaunchedInstance)

	if err != nil {
		return nil, shared.NewGeneralError("CreateInstance", err)
	}

	return launchedInstance, nil
}

func (p PublicCloudRepository) UpdateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, *shared.RepositoryError) {

	updateInstanceOpts, err := p.convertEntityToUpdateInstanceOpts(instance)
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

	updatedInstance, err := p.convertInstanceDetails(*sdkUpdatedInstance)
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
	response, err := p.publicCLoudAPI.TerminateInstance(p.authContext(ctx), id.String()).Execute()
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

	sdkInstanceTypes, response, err := p.publicCLoudAPI.GetUpdateInstanceTypeList(p.authContext(ctx), id.String()).Execute()
	if err != nil {
		return nil, shared.NewSdkError(
			fmt.Sprintf("GetAvailableInstanceTypesForUpdate %q", id),
			err,
			response,
		)
	}

	for _, sdkInstanceType := range sdkInstanceTypes.InstanceTypes {
		instanceTypes = append(
			instanceTypes,
			p.convertInstanceType(sdkInstanceType),
		)
	}

	return instanceTypes, nil
}

func (p PublicCloudRepository) GetRegions(ctx context.Context) (
	domain.Regions,
	*shared.RepositoryError,
) {
	var regions domain.Regions

	sdkRegions, response, err := p.publicCLoudAPI.GetRegionList(p.authContext(ctx)).Execute()
	if err != nil {
		return nil, shared.NewSdkError(
			"GetRegions",
			err,
			response,
		)
	}

	for _, sdkRegion := range sdkRegions.Regions {
		regions = append(regions, p.convertRegion(sdkRegion))
	}

	return regions, nil
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
		publicCLoudAPI:                    client.PublicCloudAPI,
		token:                             token,
		convertInstanceDetails:            convertInstanceDetails,
		convertInstance:                   convertInstance,
		convertAutoScalingGroupDetails:    convertAutoScalingGroupDetails,
		convertLoadBalancerDetails:        convertLoadBalancerDetails,
		convertInstanceType:               convertInstanceType,
		convertRegion:                     convertRegion,
		convertEntityToLaunchInstanceOpts: convertEntityToLaunchInstanceOpts,
		convertEntityToUpdateInstanceOpts: convertEntityToUpdateInstanceOpts,
	}
}
