# Code architecture

## Packages

Each package corresponds to a product group, i.e.: `publiccloud`

## DataSources

### DataSource files

All code pertaining to dataSources, including models, belongs in the dataSource
file.
The format for this file is `<ENDPOINT>_data_source.go`.
For instances the filename would be `instances_data_source.go`

### DataSource structs

DataSource structs should adhere to the following naming convention:
`<NAME>DataSource`.
For `instances` this would be `instancesDataSource`

### DataSource models

DataSource model structs should adhere to the following convention:
`<MODEL_NAME>DataSourceModel`.
The `instances` data model would be named `instancesDataSourceModel`.

#### DataSource adaptation functions

Adapt functions to convert SDK models to dataSource models have the following
naming convention: `adapt<SDK_MODEL_NAME>To<DATASOURCE_MODEL_NAME>DataSource`.
The function to adapt an SDK Instance to an Instance Datasource would be
named `adaptInstanceToInstanceDataSource`.

## Resources

### Resource files

All code pertaining to resources, including models, belongs in the resource
file.
The format for this file is `<ENDPOINT>_resource.go`.
For instances the filename would be `instance_resource.go`

### Resource structs

Resource structs should adhere to the following naming convention:
`<NAME>Resource`.
For `instance` this would be `instanceResource`.

### Resource models

Resource model structs should adhere to the following convention:
`<MODEL_NAME>ResourceModel`.
The `instance` data model would be named `instanceResourceModel`

#### Resource adaptation functions

Adapt functions to convert SDK models to resource models have the following
naming convention: `adapt<SDK_MODEL_NAME>To<RESOURCE_MODEL_NAME>Resource`.
The function to adapt an SDK Instance to an Instance Resource would be
named `adaptInstanceToInstanceResource`.

## Validators

As validators are often shared between resource they belong in the `validators.go`
file.

## SDK

Where possible, use the SDK getters.
Instead of `sdkInstance.Id` use `sdkInstance.GetId()`.

## Metadata

For both resources & data sources,
the following naming conventions apply for type names:

- `response.TypeName` names must be lowercase
- an underscore must replace non-alphabetic characters

`Public Cloud Load Balancers` thus becomes `public_cloud_load_balancers`

## Configure

If possible, map the SDK API as the client used by Terraform.

```go
func (i *imageResource) Configure(
	_ context.Context,
	request resource.ConfigureRequest,
	response *resource.ConfigureResponse,
) {
  ...

	coreClient, ok := request.ProviderData.(client.Client)

  ...

	i.client = coreClient.PublicCloudAPI
}
```
