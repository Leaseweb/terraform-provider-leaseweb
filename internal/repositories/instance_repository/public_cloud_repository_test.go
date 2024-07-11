package instance_repository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/domain/entity"
	"terraform-provider-leaseweb/internal/core/shared/value_object"
)

var (
	_ publicCloudApi = &publicCloudApiSpy{}
)

type publicCloudApiSpy struct {
	instance                        *publicCloud.InstanceDetails
	updatedInstance                 *publicCloud.InstanceDetails
	instanceList                    []publicCloud.Instance
	autoScalingGroup                *publicCloud.AutoScalingGroupDetails
	loadBalancer                    *publicCloud.LoadBalancerDetails
	launchedInstance                *publicCloud.Instance
	getInstanceListExecuteError     error
	getInstanceExecuteError         error
	getAutoScalingGroupExecuteError error
	getLoadBalancerExecuteError     error
	launchInstanceExecuteError      error
	updateInstanceExecuteError      error
	terminateInstanceExecuteError   error
	wantedGetInstanceId             *string
	wantedGetAutoScalingGroupId     *string
	wantedGetLoadBalancerId         *string
	wantedTerminateInstanceId       *string
	t                               *testing.T
}

func (p publicCloudApiSpy) TerminateInstance(
	ctx context.Context,
	instanceId string,
) publicCloud.ApiTerminateInstanceRequest {
	if p.wantedTerminateInstanceId != nil {
		assert.Equal(p.t, *p.wantedTerminateInstanceId, instanceId)
	}

	return publicCloud.ApiTerminateInstanceRequest{}
}

func (p publicCloudApiSpy) TerminateInstanceExecute(r publicCloud.ApiTerminateInstanceRequest) (
	*http.Response,
	error,
) {
	return nil, p.terminateInstanceExecuteError
}

func (p publicCloudApiSpy) LaunchInstance(ctx context.Context) publicCloud.ApiLaunchInstanceRequest {
	return publicCloud.ApiLaunchInstanceRequest{}
}

func (p publicCloudApiSpy) LaunchInstanceExecute(r publicCloud.ApiLaunchInstanceRequest) (
	*publicCloud.Instance,
	*http.Response,
	error,
) {
	return p.launchedInstance, nil, p.launchInstanceExecuteError
}

func (p publicCloudApiSpy) UpdateInstance(
	ctx context.Context,
	instanceId string,
) publicCloud.ApiUpdateInstanceRequest {
	return publicCloud.ApiUpdateInstanceRequest{}
}

func (p publicCloudApiSpy) UpdateInstanceExecute(r publicCloud.ApiUpdateInstanceRequest) (
	*publicCloud.InstanceDetails,
	*http.Response,
	error,
) {
	return p.updatedInstance, nil, p.updateInstanceExecuteError
}

func (p publicCloudApiSpy) GetInstanceList(ctx context.Context) publicCloud.ApiGetInstanceListRequest {
	return publicCloud.ApiGetInstanceListRequest{}
}

func (p publicCloudApiSpy) GetInstanceListExecute(r publicCloud.ApiGetInstanceListRequest) (
	*publicCloud.GetInstanceListResult,
	*http.Response,
	error,
) {
	return &publicCloud.GetInstanceListResult{Instances: p.instanceList},
		nil,
		p.getInstanceListExecuteError
}

func (p publicCloudApiSpy) GetAutoScalingGroup(
	ctx context.Context,
	autoScalingGroupId string,
) publicCloud.ApiGetAutoScalingGroupRequest {
	if p.wantedGetAutoScalingGroupId != nil {
		assert.Equal(p.t, *p.wantedGetAutoScalingGroupId, autoScalingGroupId)
	}

	return publicCloud.ApiGetAutoScalingGroupRequest{}
}

func (p publicCloudApiSpy) GetAutoScalingGroupExecute(r publicCloud.ApiGetAutoScalingGroupRequest) (
	*publicCloud.AutoScalingGroupDetails,
	*http.Response,
	error,
) {

	return p.autoScalingGroup, nil, p.getAutoScalingGroupExecuteError
}

func (p publicCloudApiSpy) GetLoadBalancer(
	ctx context.Context,
	loadBalancerId string,
) publicCloud.ApiGetLoadBalancerRequest {
	if p.wantedGetLoadBalancerId != nil {
		assert.Equal(p.t, *p.wantedGetLoadBalancerId, loadBalancerId)
	}

	return publicCloud.ApiGetLoadBalancerRequest{}
}

func (p publicCloudApiSpy) GetLoadBalancerExecute(r publicCloud.ApiGetLoadBalancerRequest) (
	*publicCloud.LoadBalancerDetails,
	*http.Response,
	error,
) {
	return p.loadBalancer, nil, p.getLoadBalancerExecuteError
}

func (p publicCloudApiSpy) GetInstance(
	ctx context.Context,
	instanceId string,
) publicCloud.ApiGetInstanceRequest {
	if p.wantedGetInstanceId != nil {
		assert.Equal(p.t, *p.wantedGetInstanceId, instanceId)
	}

	return publicCloud.ApiGetInstanceRequest{}
}

func (p publicCloudApiSpy) GetInstanceExecute(r publicCloud.ApiGetInstanceRequest) (
	*publicCloud.InstanceDetails,
	*http.Response,
	error,
) {
	return p.instance, nil, p.getInstanceExecuteError
}

func TestNewPublicCloudRepository(t *testing.T) {
	t.Run("token is set properly", func(t *testing.T) {
		got := NewPublicCloudRepository("token", Optional{})

		assert.Equal(t, "token", got.token)
	})
}

func TestPublicCloudRepository_authContext(t *testing.T) {
	publicCloudRepository := NewPublicCloudRepository("token", Optional{})
	got := publicCloudRepository.authContext(context.TODO()).Value(
		publicCloud.ContextAPIKeys,
	)

	assert.Equal(
		t,
		map[string]publicCloud.APIKey{"X-LSW-Auth": {Key: "token", Prefix: ""}},
		got,
	)
}

func TestPublicCloudRepository_GetInstance(t *testing.T) {
	t.Run("expected instance entity is returned", func(t *testing.T) {
		id := value_object.NewGeneratedUuid()
		expectedId := id.String()

		sdkInstance := publicCloud.InstanceDetails{Id: value_object.NewGeneratedUuid().String()}
		expected := entity.Instance{}

		apiSpy := publicCloudApiSpy{
			instance:            &sdkInstance,
			wantedGetInstanceId: &expectedId,
			t:                   t,
		}

		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
			convertInstance: func(
				sdkInstance publicCloud.InstanceDetails,
				autoScalingGroup *entity.AutoScalingGroup,
			) (*entity.Instance, error) {
				assert.Equal(
					t,
					apiSpy.instance,
					&sdkInstance,
					"sdkInstance is converted",
				)

				return &expected, nil
			},
		}

		got, err := publicCloudRepository.GetInstance(id, context.TODO())

		assert.NoError(t, err)
		assert.Same(t, &expected, got)
	})

	t.Run(
		"error is returned if instance cannot be retrieved from the sdk",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				getInstanceExecuteError: errors.New("error getting instance"),
			}

			PublicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := PublicCloudRepository.GetInstance(value_object.NewGeneratedUuid(), context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting instance")
		},
	)

	t.Run("expected autoScalingGroup is set", func(t *testing.T) {
		autoScalingGroupId := value_object.NewGeneratedUuid()
		convertedAutoScalingGroupId := value_object.NewGeneratedUuid()

		apiSpy := publicCloudApiSpy{
			instance: &publicCloud.InstanceDetails{
				AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(&publicCloud.AutoScalingGroup{
					Id: autoScalingGroupId.String()},
				),
			},
			autoScalingGroup: &publicCloud.AutoScalingGroupDetails{Id: autoScalingGroupId.String()},
		}

		PublicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
			convertInstance: func(
				skInstance publicCloud.InstanceDetails,
				autoScalingGroup *entity.AutoScalingGroup,
			) (*entity.Instance, error) {
				assert.Equal(
					t,
					convertedAutoScalingGroupId,
					autoScalingGroup.Id,
					"autoScalingGroup is passed on to convertInstance",
				)

				return &entity.Instance{AutoScalingGroup: &entity.AutoScalingGroup{Id: convertedAutoScalingGroupId}}, nil
			},
			convertAutoScalingGroup: func(
				sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
				loadBalancer *entity.LoadBalancer,
			) (*entity.AutoScalingGroup, error) {
				assert.Equal(
					t,
					autoScalingGroupId.String(),
					sdkAutoScalingGroup.GetId(),
					"sdkAutoScalingGroup is converted",
				)
				return &entity.AutoScalingGroup{Id: convertedAutoScalingGroupId}, nil
			},
		}

		got, err := PublicCloudRepository.GetInstance(value_object.NewGeneratedUuid(), context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, convertedAutoScalingGroupId, got.AutoScalingGroup.Id)
	})

	t.Run(
		"error is returned if autoScalingGroup uuid is invalid",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				instance: &publicCloud.InstanceDetails{
					AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(
						&publicCloud.AutoScalingGroup{
							Id: "tralala",
						},
					),
				},
			}

			PublicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := PublicCloudRepository.GetInstance(value_object.NewGeneratedUuid(), context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "cannot convert string to uuid")
		},
	)

	t.Run(
		"error is returned if autoScalingGroup cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				instance: &publicCloud.InstanceDetails{
					AutoScalingGroup: *publicCloud.NewNullableAutoScalingGroup(
						&publicCloud.AutoScalingGroup{
							Id: value_object.NewGeneratedUuid().String(),
						},
					),
				},
				getAutoScalingGroupExecuteError: errors.New("error getting autoScalingGroup"),
			}

			PublicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := PublicCloudRepository.GetInstance(value_object.NewGeneratedUuid(), context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting autoScalingGroup")
		},
	)
}

func TestPublicCLoudRepository_GetLoadBalancer(t *testing.T) {
	t.Run("expected loadBalancer entity is returned", func(t *testing.T) {
		id := value_object.NewGeneratedUuid()
		wantedLoadBalancerId := id.String()

		convertedId := value_object.NewGeneratedUuid()

		apiSpy := publicCloudApiSpy{
			loadBalancer:            &publicCloud.LoadBalancerDetails{Id: id.String()},
			wantedGetLoadBalancerId: &wantedLoadBalancerId,
			t:                       t,
		}

		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
			convertLoadBalancer: func(
				sdkLoadBalancer publicCloud.LoadBalancerDetails,
			) (*entity.LoadBalancer, error) {
				assert.Equal(
					t,
					id.String(),
					sdkLoadBalancer.Id,
					"sdkLoadBalancer is passed on to convertLoadBalancer",
				)

				return &entity.LoadBalancer{Id: convertedId}, nil
			},
		}

		got, err := publicCloudRepository.GetLoadBalancer(id, context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, convertedId, got.Id)
	})

	t.Run(
		"error is returned when loadBalancer cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				getLoadBalancerExecuteError: errors.New("error getting loadBalancer"),
			}

			PublicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := PublicCloudRepository.GetLoadBalancer(value_object.NewGeneratedUuid(), context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting loadBalancer")
		},
	)

	t.Run(
		"error is returned if loadBalancer cannot be converted",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				loadBalancer: &publicCloud.LoadBalancerDetails{Id: value_object.NewGeneratedUuid().String()},
			}

			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertLoadBalancer: func(
					sdkLoadBalancer publicCloud.LoadBalancerDetails,
				) (*entity.LoadBalancer, error) {

					return nil, errors.New("conversion error")
				},
			}

			_, err := publicCloudRepository.GetLoadBalancer(
				value_object.NewGeneratedUuid(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "conversion error")
		},
	)
}

func TestPublicCloudRepository_GetAutoScalingGroup(t *testing.T) {
	t.Run(
		"expected autoScalingGroup entity is returned",
		func(t *testing.T) {
			id := value_object.NewGeneratedUuid()
			wantedAutoScalingGroupId := id.String()
			convertedId := value_object.NewGeneratedUuid()

			apiSpy := publicCloudApiSpy{
				autoScalingGroup:            &publicCloud.AutoScalingGroupDetails{Id: id.String()},
				wantedGetAutoScalingGroupId: &wantedAutoScalingGroupId,
				t:                           t,
			}

			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertAutoScalingGroup: func(
					sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
					loadBalancer *entity.LoadBalancer,
				) (*entity.AutoScalingGroup, error) {
					assert.Equal(
						t,
						id.String(),
						sdkAutoScalingGroup.Id,
						"sdkLoadBalancer is passed on to convertLoadBalancer",
					)

					return &entity.AutoScalingGroup{Id: convertedId}, nil
				},
			}

			got, err := publicCloudRepository.GetAutoScalingGroup(id, context.TODO())

			assert.NoError(t, err)
			assert.Equal(t, convertedId, got.Id)
		},
	)

	t.Run(
		"return error if autoScalingGroup cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				getAutoScalingGroupExecuteError: errors.New("error getting autoScalingGroup"),
			}

			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := publicCloudRepository.GetAutoScalingGroup(
				value_object.NewGeneratedUuid(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting autoScalingGroup")
		},
	)

	t.Run("return error if loadBalancer id is invalid", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{
			autoScalingGroup: &publicCloud.AutoScalingGroupDetails{
				LoadBalancer: *publicCloud.NewNullableLoadBalancer(
					&publicCloud.LoadBalancer{Id: "tralala"},
				),
			},
		}

		publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

		_, err := publicCloudRepository.GetAutoScalingGroup(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot convert string to uuid")
	},
	)

	t.Run(
		"return error if loadBalancer cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				autoScalingGroup: &publicCloud.AutoScalingGroupDetails{
					LoadBalancer: *publicCloud.NewNullableLoadBalancer(
						&publicCloud.LoadBalancer{
							Id: value_object.NewGeneratedUuid().String(),
						},
					),
				},
				getLoadBalancerExecuteError: errors.New("error getting loadBalancer"),
			}

			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := publicCloudRepository.GetAutoScalingGroup(
				value_object.NewGeneratedUuid(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting loadBalancer")
		},
	)

	t.Run(
		"return error if autoScalingGroup cannot be converted",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				autoScalingGroup: &publicCloud.AutoScalingGroupDetails{},
			}

			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertAutoScalingGroup: func(
					sdkAutoScalingGroup publicCloud.AutoScalingGroupDetails,
					loadBalancer *entity.LoadBalancer,
				) (*entity.AutoScalingGroup, error) {
					return nil, errors.New("conversion error")
				},
			}

			_, err := publicCloudRepository.GetAutoScalingGroup(
				value_object.NewGeneratedUuid(),
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "conversion error")
		},
	)

}

func TestPublicCloudRepository_GetAllInstances(t *testing.T) {
	t.Run("expected instances entity is returned", func(t *testing.T) {
		id := value_object.NewGeneratedUuid()
		wantedInstanceId := id.String()

		apiSpy := publicCloudApiSpy{
			instanceList:        []publicCloud.Instance{{Id: id.String()}},
			instance:            &publicCloud.InstanceDetails{Id: id.String()},
			wantedGetInstanceId: &wantedInstanceId,
			t:                   t,
		}
		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
			convertInstance: func(
				sdkInstance publicCloud.InstanceDetails,
				sdkAutoScalingGroup *entity.AutoScalingGroup,
			) (*entity.Instance, error) {
				return &entity.Instance{Id: id}, nil
			},
		}

		got, err := publicCloudRepository.GetAllInstances(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, id, got[0].Id)
	})

	t.Run(
		"return error when instances cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				getInstanceListExecuteError: errors.New("error getting instances"),
			}
			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := publicCloudRepository.GetAllInstances(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting instances")
		},
	)

	t.Run(
		"return error when instance id cannot be parsed",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				instanceList: []publicCloud.Instance{{Id: "tralala"}},
			}
			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := publicCloudRepository.GetAllInstances(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"return error when instance cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				instanceList: []publicCloud.Instance{
					{Id: value_object.NewGeneratedUuid().String()},
				},
				getInstanceExecuteError: errors.New("error getting instance"),
			}
			publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}

			_, err := publicCloudRepository.GetAllInstances(context.TODO())

			assert.Error(t, err)
			assert.ErrorContains(t, err, "error getting instance")
		},
	)
}

func TestPublicCloudRepository_CreateInstance(t *testing.T) {
	t.Run("expected instance entity is created", func(t *testing.T) {
		passedInstance := entity.Instance{Id: value_object.NewGeneratedUuid()}

		expected := entity.Instance{}

		apiSpy := publicCloudApiSpy{
			launchedInstance: &publicCloud.Instance{Id: value_object.NewGeneratedUuid().String()},
			instance:         &publicCloud.InstanceDetails{},
		}
		publicCloudRepository := PublicCloudRepository{
			publicCLoudAPI: apiSpy,
			convertInstance: func(
				sdkInstance publicCloud.InstanceDetails,
				sdkAutoScalingGroup *entity.AutoScalingGroup,
			) (*entity.Instance, error) {
				return &expected, nil
			},
			convertEntityToLaunchInstanceOpts: func(instance entity.Instance) (
				*publicCloud.LaunchInstanceOpts,
				error,
			) {
				assert.Equal(
					t,
					passedInstance,
					instance,
					"passed instance entity is converted to launchInstanceOpts",
				)
				return &publicCloud.LaunchInstanceOpts{}, nil
			},
		}

		got, err := publicCloudRepository.CreateInstance(
			passedInstance,
			context.TODO(),
		)

		assert.NoError(t, err)
		assert.Same(t, &expected, got)
	})

	t.Run(
		"error is returned when launchInstanceOpts cannot be created",
		func(t *testing.T) {
			publicCloudRepository := PublicCloudRepository{
				convertEntityToLaunchInstanceOpts: func(instance entity.Instance) (
					*publicCloud.LaunchInstanceOpts,
					error,
				) {
					return nil, errors.New("error getting launchInstanceOpts")
				},
			}

			_, err := publicCloudRepository.CreateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.ErrorContains(t, err, "error getting launchInstanceOpts")
		},
	)

	t.Run(
		"error is returned when instance cannot be launched in sdk",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				launchInstanceExecuteError: errors.New("some error"),
			}
			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertEntityToLaunchInstanceOpts: func(instance entity.Instance) (
					*publicCloud.LaunchInstanceOpts,
					error,
				) {
					return &publicCloud.LaunchInstanceOpts{}, nil
				},
			}

			_, err := publicCloudRepository.CreateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error is returned when id of launched instance is invalid",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				launchedInstance: &publicCloud.Instance{Id: "tralala"},
			}
			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertEntityToLaunchInstanceOpts: func(instance entity.Instance) (
					*publicCloud.LaunchInstanceOpts,
					error,
				) {
					return &publicCloud.LaunchInstanceOpts{}, nil
				},
			}

			_, err := publicCloudRepository.CreateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"error is returned when instanceDetails cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				launchedInstance: &publicCloud.Instance{
					Id: value_object.NewGeneratedUuid().String(),
				},
				getInstanceExecuteError: errors.New("some error"),
			}
			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertEntityToLaunchInstanceOpts: func(instance entity.Instance) (
					*publicCloud.LaunchInstanceOpts,
					error,
				) {
					return &publicCloud.LaunchInstanceOpts{}, nil
				},
			}

			_, err := publicCloudRepository.CreateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)
}

func TestPublicCloudRepository_UpdateInstance(t *testing.T) {
	t.Run(
		"expected instance entity is returned on update",
		func(t *testing.T) {
			passedInstance := entity.Instance{Id: value_object.NewGeneratedUuid()}

			expected := entity.Instance{}

			apiSpy := publicCloudApiSpy{
				updatedInstance: &publicCloud.InstanceDetails{
					Id: value_object.NewGeneratedUuid().String(),
				},
				instance: &publicCloud.InstanceDetails{},
			}
			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertInstance: func(
					sdkInstance publicCloud.InstanceDetails,
					sdkAutoScalingGroup *entity.AutoScalingGroup,
				) (*entity.Instance, error) {
					return &expected, nil
				},
				convertEntityToUpdateInstanceOpts: func(instance entity.Instance) (
					*publicCloud.UpdateInstanceOpts,
					error,
				) {
					assert.Equal(
						t,
						passedInstance,
						instance,
						"passed instance entity is converted to updateInstanceOpts",
					)
					return &publicCloud.UpdateInstanceOpts{}, nil
				},
			}

			got, err := publicCloudRepository.UpdateInstance(
				passedInstance,
				context.TODO(),
			)

			assert.NoError(t, err)
			assert.Same(t, &expected, got)
		},
	)

	t.Run(
		"error is returned when updateInstanceOpts cannot be created",
		func(t *testing.T) {
			publicCloudRepository := PublicCloudRepository{
				convertEntityToUpdateInstanceOpts: func(instance entity.Instance) (
					*publicCloud.UpdateInstanceOpts,
					error,
				) {
					return nil, errors.New("error getting updateInstanceOpts")
				},
			}

			_, err := publicCloudRepository.UpdateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.ErrorContains(t, err, "error getting updateInstanceOpts")
		},
	)

	t.Run(
		"error is returned when instance cannot be updated in sdk",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				updateInstanceExecuteError: errors.New("some error"),
			}
			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertEntityToUpdateInstanceOpts: func(instance entity.Instance) (
					*publicCloud.UpdateInstanceOpts,
					error,
				) {
					return &publicCloud.UpdateInstanceOpts{}, nil
				},
			}

			_, err := publicCloudRepository.UpdateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

	t.Run(
		"error is returned when id of updated instance is invalid",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				updatedInstance: &publicCloud.InstanceDetails{Id: "tralala"},
			}
			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertEntityToUpdateInstanceOpts: func(instance entity.Instance) (
					*publicCloud.UpdateInstanceOpts,
					error,
				) {
					return &publicCloud.UpdateInstanceOpts{}, nil
				},
			}

			_, err := publicCloudRepository.UpdateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "tralala")
		},
	)

	t.Run(
		"error is returned when instanceDetails cannot be retrieved",
		func(t *testing.T) {
			apiSpy := publicCloudApiSpy{
				updatedInstance: &publicCloud.InstanceDetails{
					Id: value_object.NewGeneratedUuid().String(),
				},
				getInstanceExecuteError: errors.New("some error"),
			}
			publicCloudRepository := PublicCloudRepository{
				publicCLoudAPI: apiSpy,
				convertEntityToUpdateInstanceOpts: func(instance entity.Instance) (
					*publicCloud.UpdateInstanceOpts,
					error,
				) {
					return &publicCloud.UpdateInstanceOpts{}, nil
				},
			}

			_, err := publicCloudRepository.UpdateInstance(
				entity.Instance{},
				context.TODO(),
			)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "some error")
		},
	)

}

func TestPublicCloudRepository_DeleteInstance(t *testing.T) {
	t.Run("expected instance entity is deleted", func(t *testing.T) {
		id := value_object.NewGeneratedUuid()
		wantedId := id.String()

		apiSpy := publicCloudApiSpy{wantedTerminateInstanceId: &wantedId, t: t}

		publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}
		err := publicCloudRepository.DeleteInstance(id, context.TODO())

		assert.NoError(t, err)
	})

	t.Run("error is returned when instance deletion fails", func(t *testing.T) {
		apiSpy := publicCloudApiSpy{terminateInstanceExecuteError: errors.New("some error")}

		publicCloudRepository := PublicCloudRepository{publicCLoudAPI: apiSpy}
		err := publicCloudRepository.DeleteInstance(
			value_object.NewGeneratedUuid(),
			context.TODO(),
		)

		assert.ErrorIs(t, ErrSomethingWentWrongDeletingTheInstance, err)

	})
}
