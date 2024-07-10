package instance_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/core/domain/entity"
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
	request := p.publicCLoudAPI.GetInstanceList(p.authContext(ctx))
	result, _, err := p.publicCLoudAPI.GetInstanceListExecute(request)

	if err != nil {
		return nil, fmt.Errorf("cannot retrieve instances: %w", err)
	}

	var instances entity.Instances

	for _, instance := range result.Instances {
		instanceId, err := uuid.Parse(instance.GetId())
		if err != nil {
			return nil, fmt.Errorf(
				"cannot parse uuid %s: %w",
				instance.GetId(),
				err,
			)

		}
		instanceDetail, err := p.GetInstance(instanceId, ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot retrieve instance details %s: %w",
				instanceId,
				err,
			)
		}

		instances = append(instances, *instanceDetail)
	}

	return instances, nil
}

func (p PublicCloudRepository) GetInstance(
	id uuid.UUID,
	ctx context.Context,
) (*entity.Instance, error) {
	var autoScalingGroup *entity.AutoScalingGroup

	request := p.publicCLoudAPI.GetInstance(p.authContext(ctx), id.String())
	instance, _, err := p.publicCLoudAPI.GetInstanceExecute(request)

	if err != nil {
		return nil, fmt.Errorf("cannot retrieve instance: %w", err)
	}

	sdkAutoScalingGroup, _ := instance.GetAutoScalingGroupOk()

	if sdkAutoScalingGroup != nil {
		autoScalingGroupId, err := convertStringToUuid(sdkAutoScalingGroup.GetId())
		if err != nil {
			return nil, fmt.Errorf(
				"error parsing autoScalingGroup id %q: %w",
				sdkAutoScalingGroup.GetId(),
				err,
			)
		}
		autoScalingGroup, err = p.GetAutoScalingGroup(
			*autoScalingGroupId,
			ctx,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot retrieve auto scaling group details %s: %w",
				sdkAutoScalingGroup.GetId(),
				err,
			)
		}
	}

	return p.convertInstance(*instance, autoScalingGroup)
}

func (p PublicCloudRepository) GetAutoScalingGroup(
	id uuid.UUID,
	ctx context.Context,
) (*entity.AutoScalingGroup, error) {
	var loadBalancer *entity.LoadBalancer

	request := p.publicCLoudAPI.GetAutoScalingGroup(
		p.authContext(ctx),
		id.String(),
	)

	sdkAutoScalingGroup, _, err := p.publicCLoudAPI.GetAutoScalingGroupExecute(request)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot retrieve auto scaling group details %s: %w",
			id,
			err,
		)
	}

	if sdkAutoScalingGroup.LoadBalancer.Get() != nil {
		loadBalancerId, err := convertStringToUuid(sdkAutoScalingGroup.LoadBalancer.Get().GetId())
		if err != nil {
			return nil, fmt.Errorf(
				"error parsing sdkAutoScalingGroup id %q: %w",
				sdkAutoScalingGroup.LoadBalancer.Get().GetId(),
				err,
			)
		}

		loadBalancer, err = p.GetLoadBalancer(*loadBalancerId, ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot retrieve auto loadBalancer details %s: %w",
				sdkAutoScalingGroup.LoadBalancer.Get().GetId(),
				err,
			)
		}
	}

	autoScalingGroupEntity, err := p.convertAutoScalingGroup(
		*sdkAutoScalingGroup,
		loadBalancer,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot convert sdkAutoScalingGroup %s: %w",
			id,
			err,
		)
	}

	return autoScalingGroupEntity, nil
}

func (p PublicCloudRepository) GetLoadBalancer(
	id uuid.UUID,
	ctx context.Context,
) (*entity.LoadBalancer, error) {
	var loadBalancer *entity.LoadBalancer

	request := p.publicCLoudAPI.GetLoadBalancer(p.authContext(ctx), id.String())

	sdkLoadBalancer, _, err := p.publicCLoudAPI.GetLoadBalancerExecute(request)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot retrieve auto loadBalancer details %s: %w",
			id,
			err,
		)
	}

	loadBalancer, err = p.convertLoadBalancer(*sdkLoadBalancer)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot convert loadBalancer %s: %w",
			id,
			err,
		)
	}

	return loadBalancer, nil
}

func NewPublicCloudRepository(token string, optional Optional) PublicCloudRepository {
	configuration := publicCloud.NewConfiguration()

	if optional.Host != nil {
		configuration.Host = *optional.Host
	}
	if optional.Scheme != nil {
		configuration.Scheme = *optional.Scheme
	}

	client := *publicCloud.NewAPIClient(configuration)

	return PublicCloudRepository{
		publicCLoudAPI:          client.PublicCloudAPI,
		token:                   token,
		convertInstance:         convertInstance,
		convertAutoScalingGroup: convertAutoScalingGroup,
		convertLoadBalancer:     convertLoadBalancer,
	}
}
