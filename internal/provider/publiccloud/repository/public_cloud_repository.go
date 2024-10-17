// Package repository implements repository logic to access the public_cloud sdk.
package repository

import (
	"context"
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	sharedRepository "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/repository"
)

// Optional contains optional values that can be passed to NewPublicCloudRepository.
type Optional struct {
	Host   *string
	Scheme *string
}

// PublicCloudRepository fulfills contract for ports.PublicCloudRepository.
type PublicCloudRepository struct {
	publicCLoudAPI      publicCloud.PublicCloudAPI
	token               string
	cachedInstanceTypes sharedRepository.SyncedMap[string, []string]
	cachedRegions       sharedRepository.SyncedMap[string, []string]
}

// Injects the authentication token into the context for the sdk.
func (p *PublicCloudRepository) authContext(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		publicCloud.ContextAPIKeys,
		map[string]publicCloud.APIKey{
			"X-LSW-Auth": {Key: p.token, Prefix: ""},
		},
	)
}

func (p *PublicCloudRepository) GetAllInstances(ctx context.Context) (
	[]publicCloud.Instance,
	*sharedRepository.RepositoryError,
) {
	var instances []publicCloud.Instance

	request := p.publicCLoudAPI.GetInstanceList(p.authContext(ctx))

	result, response, err := request.Execute()

	if err != nil {
		return nil, sharedRepository.NewSdkError("GetAllInstances", err, response)
	}

	metadata := result.GetMetadata()
	pagination := sharedRepository.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, sharedRepository.NewSdkError("GetAllInstances", err, response)
		}

		instances = append(instances, result.Instances...)

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, sharedRepository.NewSdkError("GetAllInstances", err, response)
		}
	}

	return instances, nil
}

func (p *PublicCloudRepository) GetInstance(
	id string,
	ctx context.Context,
) (*publicCloud.InstanceDetails, *sharedRepository.RepositoryError) {
	instance, response, err := p.publicCLoudAPI.GetInstance(
		p.authContext(ctx),
		id,
	).Execute()

	if err != nil {
		return nil, sharedRepository.NewSdkError(
			fmt.Sprintf("GetInstance %q", id),
			err,
			response,
		)
	}

	return instance, nil
}

func (p *PublicCloudRepository) LaunchInstance(
	opts publicCloud.LaunchInstanceOpts,
	ctx context.Context,
) (*publicCloud.Instance, *sharedRepository.RepositoryError) {
	instance, response, err := p.publicCLoudAPI.
		LaunchInstance(p.authContext(ctx)).
		LaunchInstanceOpts(opts).Execute()

	if err != nil {
		return nil, sharedRepository.NewSdkError(
			"LaunchInstance",
			err,
			response,
		)
	}

	return instance, nil
}

func (p *PublicCloudRepository) UpdateInstance(
	id string,
	opts publicCloud.UpdateInstanceOpts,
	ctx context.Context,
) (*publicCloud.InstanceDetails, *sharedRepository.RepositoryError) {
	instance, response, err := p.publicCLoudAPI.UpdateInstance(
		p.authContext(ctx),
		id,
	).UpdateInstanceOpts(opts).Execute()
	if err != nil {
		return nil, sharedRepository.NewSdkError(
			fmt.Sprintf("UpdateInstance %q", id),
			err,
			response,
		)
	}

	return instance, nil
}

func (p *PublicCloudRepository) DeleteInstance(
	id string,
	ctx context.Context,
) *sharedRepository.RepositoryError {
	response, err := p.publicCLoudAPI.TerminateInstance(
		p.authContext(ctx),
		id,
	).Execute()
	if err != nil {
		return sharedRepository.NewSdkError(
			fmt.Sprintf("DeleteInstance %q", id),
			err,
			response,
		)
	}

	return nil
}

func (p *PublicCloudRepository) GetAvailableInstanceTypesForUpdate(
	id string,
	ctx context.Context,
) ([]string, *sharedRepository.RepositoryError) {
	var instanceTypes []string

	sdkInstanceTypes, response, err := p.publicCLoudAPI.GetUpdateInstanceTypeList(
		p.authContext(ctx),
		id,
	).Execute()
	if err != nil {
		return nil, sharedRepository.NewSdkError(
			fmt.Sprintf("GetAvailableInstanceTypesForUpdate %q", id),
			err,
			response,
		)
	}

	for _, sdkInstanceType := range sdkInstanceTypes.InstanceTypes {
		instanceTypes = append(instanceTypes, string(sdkInstanceType.Name))
	}

	return instanceTypes, nil
}

func (p *PublicCloudRepository) GetRegions(ctx context.Context) (
	[]string,
	*sharedRepository.RepositoryError,
) {
	var regions []string

	regions, ok := p.cachedRegions.Get("all")
	if ok {
		return regions, nil
	}

	request := p.publicCLoudAPI.GetRegionList(p.authContext(ctx))

	result, response, err := request.Execute()

	if err != nil {
		return nil, sharedRepository.NewSdkError("getRegions", err, response)
	}

	metadata := result.GetMetadata()
	pagination := sharedRepository.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, sharedRepository.NewSdkError("getRegions", err, response)
		}

		for _, sdkRegion := range result.Regions {
			regions = append(regions, string(sdkRegion.Name))
		}

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, sharedRepository.NewSdkError("GetAllInstances", err, response)
		}
	}

	p.cachedRegions.Set("all", regions)

	return regions, nil
}

func (p *PublicCloudRepository) GetInstanceTypesForRegion(
	region string,
	ctx context.Context,
) ([]string, *sharedRepository.RepositoryError) {
	var instanceTypes []string

	cachedInstanceTypes, ok := p.cachedInstanceTypes.Get(region)
	if ok {
		return cachedInstanceTypes, nil
	}

	request := p.publicCLoudAPI.GetInstanceTypeList(p.authContext(ctx)).
		Region(publicCloud.RegionName(region))

	result, response, err := request.Execute()

	if err != nil {
		return nil, sharedRepository.NewSdkError(
			"GetInstanceTypesForRegion",
			err,
			response,
		)
	}

	metadata := result.GetMetadata()
	pagination := sharedRepository.NewPagination(
		metadata.GetLimit(),
		metadata.GetTotalCount(),
		request,
	)

	for {
		result, response, err := request.Execute()
		if err != nil {
			return nil, sharedRepository.NewSdkError(
				"GetInstanceTypesForRegion",
				err,
				response,
			)
		}

		for _, sdkInstanceType := range result.InstanceTypes {
			instanceTypes = append(instanceTypes, string(sdkInstanceType.Name))
		}

		if !pagination.CanIncrement() {
			break
		}

		request, err = pagination.NextPage()
		if err != nil {
			return nil, sharedRepository.NewSdkError("GetAllInstances", err, response)
		}
	}

	p.cachedInstanceTypes.Set(region, instanceTypes)

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
		publicCLoudAPI:      client.PublicCloudAPI,
		token:               token,
		cachedInstanceTypes: sharedRepository.NewSyncedMap[string, []string](),
		cachedRegions:       sharedRepository.NewSyncedMap[string, []string](),
	}
}
