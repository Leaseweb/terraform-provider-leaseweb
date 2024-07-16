package instance_repository

import (
	"context"
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

type Optional struct {
	Host   *string
	Scheme *string
}

type PublicCloudRepository struct {
	publicCLoudAPI  publicCloudApi
	token           string
	convertInstance func(
		sdkInstance publicCloud.InstanceDetails,
		sdkAutoScalingGroup *entity.AutoScalingGroup,
	) (*entity.Instance, error)
	convertAutoScalingGroup func(
		sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
		loadBalancer *entity.LoadBalancer,
	) (*entity.AutoScalingGroup, error)
	convertLoadBalancer func(
		sdkLoadBalancer publicCloud.LoadBalancerDetails,
	) (*entity.LoadBalancer, error)
	convertEntityToLaunchInstanceOpts func(instance entity.Instance) (
		*publicCloud.LaunchInstanceOpts, error)
	convertEntityToUpdateInstanceOpts func(instance entity.Instance) (
		*publicCloud.UpdateInstanceOpts, error)
	convertRegion       func(sdkRegion publicCloud.Region) entity.Region
	convertInstanceType func(sdkInstanceType publicCloud.InstanceType) entity.InstanceType
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
	entity.Instances,
	error,
) {
	result, _, err := p.publicCLoudAPI.GetInstanceList(p.authContext(ctx)).Execute()

	if err != nil {
		return nil, fmt.Errorf("GetAllInstances: %w", err)
	}

	var instances entity.Instances

	for _, instance := range result.Instances {
		instanceId, err := value_object.NewUuid(instance.GetId())
		if err != nil {
			return nil, fmt.Errorf("GetAllInstances: %w", err)

		}
		instanceDetail, err := p.GetInstance(*instanceId, ctx)
		if err != nil {
			return nil, fmt.Errorf("GetAllInstances: %w", err)
		}

		instances = append(instances, *instanceDetail)
	}

	return instances, nil
}

func (p PublicCloudRepository) GetInstance(
	id value_object.Uuid,
	ctx context.Context,
) (*entity.Instance, error) {
	var autoScalingGroup *entity.AutoScalingGroup

	instance, _, err := p.publicCLoudAPI.GetInstance(p.authContext(ctx), id.String()).Execute()

	if err != nil {
		return nil, fmt.Errorf("GetInstance %q: %w", id, err)
	}

	sdkAutoScalingGroup, _ := instance.GetAutoScalingGroupOk()

	if sdkAutoScalingGroup != nil {
		autoScalingGroupId, err := value_object.NewUuid(sdkAutoScalingGroup.GetId())
		if err != nil {
			return nil, fmt.Errorf("GetInstance: %w", err)
		}
		autoScalingGroup, err = p.GetAutoScalingGroup(*autoScalingGroupId, ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"GetInstance: %w",
				err,
			)
		}
	}

	return p.convertInstance(*instance, autoScalingGroup)
}

func (p PublicCloudRepository) GetAutoScalingGroup(
	id value_object.Uuid,
	ctx context.Context,
) (*entity.AutoScalingGroup, error) {
	var loadBalancer *entity.LoadBalancer

	sdkAutoScalingGroup, _, err := p.publicCLoudAPI.GetAutoScalingGroup(
		p.authContext(ctx),
		id.String(),
	).Execute()
	if err != nil {
		return nil, fmt.Errorf("GetAutoScalingGroup %q: %w", id, err)
	}

	if sdkAutoScalingGroup.LoadBalancer.Get() != nil {
		loadBalancerId, err := value_object.NewUuid(sdkAutoScalingGroup.LoadBalancer.Get().GetId())
		if err != nil {
			return nil, fmt.Errorf("GetAutoScalingGroup: %w", err)
		}

		loadBalancer, err = p.GetLoadBalancer(*loadBalancerId, ctx)
		if err != nil {
			return nil, fmt.Errorf("GetAutoScalingGroup: %w", err)
		}
	}

	autoScalingGroupEntity, err := p.convertAutoScalingGroup(
		*sdkAutoScalingGroup,
		loadBalancer,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"GetAutoScalingGroup %q: %w",
			sdkAutoScalingGroup.GetId(),
			err,
		)
	}

	return autoScalingGroupEntity, nil
}

func (p PublicCloudRepository) GetLoadBalancer(
	id value_object.Uuid,
	ctx context.Context,
) (*entity.LoadBalancer, error) {
	var loadBalancer *entity.LoadBalancer

	sdkLoadBalancer, _, err := p.publicCLoudAPI.GetLoadBalancer(p.authContext(ctx), id.String()).Execute()
	if err != nil {
		return nil, fmt.Errorf("GetLoadBalancer %q: %w", id, err)
	}

	loadBalancer, err = p.convertLoadBalancer(*sdkLoadBalancer)
	if err != nil {
		return nil, fmt.Errorf(
			"GetLoadBalancer %q: %w",
			sdkLoadBalancer.GetId(),
			err,
		)
	}

	return loadBalancer, nil
}

func (p PublicCloudRepository) CreateInstance(
	instance entity.Instance,
	ctx context.Context,
) (*entity.Instance, error) {

	launchInstanceOpts, err := p.convertEntityToLaunchInstanceOpts(instance)
	if err != nil {
		return nil, fmt.Errorf("CreateInstance : %w", err)
	}

	launchedInstance, _, err := p.publicCLoudAPI.LaunchInstance(p.authContext(ctx)).LaunchInstanceOpts(*launchInstanceOpts).Execute()

	if err != nil {
		return nil, fmt.Errorf("CreateInstance: %w", err)
	}

	instanceId, err := value_object.NewUuid(launchedInstance.GetId())
	if err != nil {
		return nil, fmt.Errorf(
			"CreateInstance: %w",
			err,
		)
	}

	instanceDetails, err := p.GetInstance(*instanceId, ctx)
	if err != nil {
		return nil, fmt.Errorf("CreateInstance: %w", err)
	}

	return instanceDetails, nil
}

func (p PublicCloudRepository) UpdateInstance(
	instance entity.Instance,
	ctx context.Context,
) (*entity.Instance, error) {

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

	instanceId, err := value_object.NewUuid(updatedInstance.GetId())
	if err != nil {
		return nil, fmt.Errorf("UpdateInstance %q: %w", instance.Id, err)
	}

	instanceDetails, err := p.GetInstance(*instanceId, ctx)
	if err != nil {
		return nil, fmt.Errorf("UpdateInstance: %w", err)
	}

	return instanceDetails, nil
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
) (entity.InstanceTypes, error) {
	var instanceTypes entity.InstanceTypes

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
	entity.Regions,
	error,
) {
	var regions entity.Regions

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
		convertInstance:                   convertInstance,
		convertAutoScalingGroup:           convertAutoScalingGroup,
		convertLoadBalancer:               convertLoadBalancer,
		convertInstanceType:               convertInstanceType,
		convertRegion:                     convertRegion,
		convertEntityToLaunchInstanceOpts: convertEntityToLaunchInstanceOpts,
		convertEntityToUpdateInstanceOpts: convertEntityToUpdateInstanceOpts,
	}
}
