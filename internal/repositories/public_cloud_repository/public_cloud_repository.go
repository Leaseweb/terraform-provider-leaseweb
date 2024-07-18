package public_cloud_repository

import (
	"context"
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/core/domain"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
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
	error,
) {
	var instances domain.Instances

	result, _, err := p.publicCLoudAPI.GetInstanceList(p.authContext(ctx)).Execute()

	if err != nil {
		return nil, fmt.Errorf("GetAllInstances: %w", err)
	}

	for _, sdkInstance := range result.Instances {
		instance, err := p.convertInstance(sdkInstance)
		if err != nil {
			return nil, fmt.Errorf("GetAllInstances: %w", err)
		}

		instances = append(instances, *instance)
	}

	return instances, nil
}

func (p PublicCloudRepository) GetInstance(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.Instance, error) {
	instanceDetails, _, err := p.publicCLoudAPI.GetInstance(
		p.authContext(ctx),
		id.String(),
	).Execute()

	if err != nil {
		return nil, fmt.Errorf("GetInstance %q: %w", id, err)
	}

	return p.convertInstanceDetails(*instanceDetails)
}

func (p PublicCloudRepository) GetAutoScalingGroup(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.AutoScalingGroup, error) {
	sdkAutoScalingGroupDetails, _, err := p.publicCLoudAPI.GetAutoScalingGroup(
		p.authContext(ctx),
		id.String(),
	).Execute()
	if err != nil {
		return nil, fmt.Errorf("GetAutoScalingGroup %q: %w", id, err)
	}

	autoScalingGroup, err := p.convertAutoScalingGroupDetails(*sdkAutoScalingGroupDetails)
	if err != nil {
		return nil, fmt.Errorf(
			"GetAutoScalingGroup %q: %w",
			sdkAutoScalingGroupDetails.GetId(),
			err,
		)
	}

	return autoScalingGroup, nil
}

func (p PublicCloudRepository) GetLoadBalancer(
	id value_object.Uuid,
	ctx context.Context,
) (*domain.LoadBalancer, error) {
	var loadBalancer *domain.LoadBalancer

	sdkLoadBalancerDetails, _, err := p.publicCLoudAPI.GetLoadBalancer(
		p.authContext(ctx),
		id.String(),
	).Execute()
	if err != nil {
		return nil, fmt.Errorf("GetLoadBalancer %q: %w", id, err)
	}

	loadBalancer, err = p.convertLoadBalancerDetails(*sdkLoadBalancerDetails)
	if err != nil {
		return nil, fmt.Errorf(
			"GetLoadBalancer %q: %w",
			sdkLoadBalancerDetails.GetId(),
			err,
		)
	}

	return loadBalancer, nil
}

func (p PublicCloudRepository) CreateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, error) {

	launchInstanceOpts, err := p.convertEntityToLaunchInstanceOpts(instance)
	if err != nil {
		return nil, fmt.Errorf("CreateInstance : %w", err)
	}

	launchedInstance, _, err := p.publicCLoudAPI.LaunchInstance(p.authContext(ctx)).LaunchInstanceOpts(*launchInstanceOpts).Execute()

	if err != nil {
		return nil, fmt.Errorf("CreateInstance: %w", err)
	}

	return p.convertInstance(*launchedInstance)
}

func (p PublicCloudRepository) UpdateInstance(
	instance domain.Instance,
	ctx context.Context,
) (*domain.Instance, error) {

	updateInstanceOpts, err := p.convertEntityToUpdateInstanceOpts(instance)
	if err != nil {
		return nil, fmt.Errorf("UpdateInstance %q: %w", instance.Id, err)
	}

	updatedInstance, _, err := p.publicCLoudAPI.UpdateInstance(
		p.authContext(ctx),
		instance.Id.String(),
	).UpdateInstanceOpts(*updateInstanceOpts).Execute()
	if err != nil {
		return nil, fmt.Errorf("UpdateInstance %q: %w", instance.Id, err)
	}

	return p.convertInstanceDetails(*updatedInstance)
}

func (p PublicCloudRepository) DeleteInstance(
	id value_object.Uuid,
	ctx context.Context,
) error {
	_, err := p.publicCLoudAPI.TerminateInstance(p.authContext(ctx), id.String()).Execute()
	if err != nil {
		return fmt.Errorf("DeleteInstance %q: %w", id, err)
	}

	return nil
}

func (p PublicCloudRepository) GetAvailableInstanceTypesForUpdate(
	id value_object.Uuid,
	ctx context.Context,
) (domain.InstanceTypes, error) {
	var instanceTypes domain.InstanceTypes

	sdkInstanceTypes, _, err := p.publicCLoudAPI.GetUpdateInstanceTypeList(p.authContext(ctx), id.String()).Execute()
	if err != nil {
		return nil, fmt.Errorf(
			"GetAvailableInstanceTypesForUpdate %q: %w",
			id,
			err,
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
	error,
) {
	var regions domain.Regions

	sdkRegions, _, err := p.publicCLoudAPI.GetRegionList(p.authContext(ctx)).Execute()
	if err != nil {
		return nil, fmt.Errorf("GetRegions: %w", err)
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
