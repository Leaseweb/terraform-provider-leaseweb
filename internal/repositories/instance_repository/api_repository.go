package instance_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/core/domain/entity"
)

type ApiRepository struct {
	client publicCloud.APIClient
	token  string
}

func (a ApiRepository) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		publicCloud.ContextAPIKeys,
		map[string]publicCloud.APIKey{
			"X-LSW-Auth": {Key: a.token, Prefix: ""},
		},
	)
}

func (a ApiRepository) GetAllInstances(ctx context.Context) (
	entity.Instances,
	error,
) {
	request := a.client.PublicCloudAPI.GetInstanceList(a.authContext(ctx))
	result, _, err := a.client.PublicCloudAPI.GetInstanceListExecute(request)

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
		instanceDetail, err := a.GetInstance(instanceId, ctx)
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

func (a ApiRepository) GetInstance(
	id uuid.UUID,
	ctx context.Context,
) (*entity.Instance, error) {
	instanceRequest := a.client.PublicCloudAPI.GetInstance(
		a.authContext(ctx),
		id.String(),
	)
	var autoScalingGroupDetails *publicCloud.AutoScalingGroupDetails
	var loadBalancerDetails *publicCloud.LoadBalancerDetails

	instance, _, err := a.client.PublicCloudAPI.GetInstanceExecute(instanceRequest)

	if err != nil {
		return nil, fmt.Errorf("cannot retrieve instance: %w", err)
	}

	autoScalingGroup, _ := instance.GetAutoScalingGroupOk()

	if autoScalingGroup != nil {
		autoScalingGroupDetailsRequest := a.client.PublicCloudAPI.GetAutoScalingGroup(
			a.authContext(ctx),
			autoScalingGroup.GetId(),
		)
		autoScalingGroupDetails, _, err = a.client.PublicCloudAPI.GetAutoScalingGroupExecute(autoScalingGroupDetailsRequest)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot retrieve auto scaling group details %s: %w",
				autoScalingGroup.GetId(),
				err,
			)
		}

		if autoScalingGroupDetails.LoadBalancer.Get() != nil {
			loadBalancerDetailsRequest := a.client.PublicCloudAPI.GetLoadBalancer(
				a.authContext(ctx),
				autoScalingGroupDetails.LoadBalancer.Get().GetId(),
			)

			loadBalancerDetails, _, err = a.client.PublicCloudAPI.GetLoadBalancerExecute(loadBalancerDetailsRequest)
			if err != nil {
				return nil, fmt.Errorf(
					"cannot retrieve auto loadBalancer details %s: %w",
					autoScalingGroupDetails.LoadBalancer.Get().GetId(),
					err,
				)
			}
		}
	}

	return convertInstance(*instance, autoScalingGroupDetails, loadBalancerDetails)
}

type Optional struct {
	Host   *string
	Scheme *string
}

func NewApiRepository(token string, optional Optional) ApiRepository {
	configuration := publicCloud.NewConfiguration()

	if optional.Host != nil {
		configuration.Host = *optional.Host
	}
	if optional.Scheme != nil {
		configuration.Scheme = *optional.Scheme
	}

	return ApiRepository{
		client: *publicCloud.NewAPIClient(configuration),
		token:  token,
	}
}
