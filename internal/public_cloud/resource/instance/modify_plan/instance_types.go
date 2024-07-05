package modify_plan

import (
	"context"
	"net/http"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"terraform-provider-leaseweb/internal/client"
)

type InstanceTypes struct {
	client client.Client
	ctx    context.Context
}

func (i InstanceTypes) GetAllowedInstanceTypes(instanceTypeId string) (
	[]string,
	*http.Response,
	error,
) {
	allowedInstanceTypesRequest := i.client.PublicCloudClient.PublicCloudAPI.
		GetUpdateInstanceTypeList(
			i.client.AuthContext(i.ctx),
			instanceTypeId,
		)
	allowedInstanceTypes, sdkResponse, err := i.client.PublicCloudClient.PublicCloudAPI.
		GetUpdateInstanceTypeListExecute(allowedInstanceTypesRequest)

	if err != nil {
		return nil, sdkResponse, err
	}

	return convertSdkInstanceTypesToString(allowedInstanceTypes.GetInstanceTypes()), nil, nil
}

func convertSdkInstanceTypesToString(sdkInstanceTypes []publicCloud.InstanceType) (instanceTypes []string) {
	for _, instanceType := range sdkInstanceTypes {
		instanceTypes = append(instanceTypes, instanceType.GetName())
	}

	return
}

func NewInstanceTypes(client client.Client, ctx context.Context) InstanceTypes {
	return InstanceTypes{
		client: client,
		ctx:    ctx,
	}
}
