# Code architecture

## Packages

Each package corresponds to a product group, i.e.: `publiccloud`

## Datasources

All code pertaining to datasources, including models, belongs in the datasource
file. The format for this file is `<PRODUCT>_<ENDPOINT>_datasource.go`. Ie: for
instances the filename would be `publiccloud_instances_datasource.go`

### Datasource models

Datasource model structs should adhere to the following convention:

`datasourceModel<MODEL_NAME>`. So the instances data model would be named
`dataSourceModelinstances`

#### Datasource adaptation functions

Adapt functions to convert SDK models to datasource models have the following
naming convention: `adaptSdk<SDKMODEL_NAME>ToDatasource<DATASOURCE_MODEL_NAME>`.
So the function to adapt an SDK Instance to an Instance Datasource would be
named `adaptSdkInstanceToDatasourceInstance`.

## Resources

All code pertaining to resources, including models, belongs in the resource
file. The format for this file is `<PRODUCT>_<ENDPOINT>_resource.go`. Ie: for
instances the filename would be `publiccloud_instances_resource.go`

### Resource models

Datasource model structs should adhere to the following convention:

`resourceModel<MODEL_NAME>`. So the instance data model would be named
`resourceModelinstance`

#### Resource adaptation functions

Adapt functions to convert SDK models to resource models have the following
naming convention: `adaptSdk<SDKMODEL_NAME>ToResource<Resource_MODEL_NAME>`.
So the function to adapt an SDK Instance to an Instance Resource would be
named `adaptSdkInstanceToResourceInstance`.

## Validators

As validators are often shared between resource they belong in the `validators.go`
file.
