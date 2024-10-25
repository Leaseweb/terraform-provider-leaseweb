# Code architecture

## Packages

Each package corresponds to a product group, i.e.: `publiccloud`

## Datasources

### Datasource files

All code pertaining to datasources, including models, belongs in the datasource
file. The format for this file is `datasource_<ENDPOINT>.go`. Ie: for
instances the filename would be `datasource_instances.go`

### Datasource structs

Structs should adhere to the following naming convention: `datasource<NAME>`.
`<NAME>` is always plural.
So for `instances` this would be `datasourceInstances`

### Datasource models

Datasource model structs should adhere to the following convention:

`datasourceModel<MODEL_NAME>`. So the `instances` data model would be named
`datasourceModelInstances`.

#### Datasource adaptation functions

Adapt functions to convert SDK models to datasource models have the following
naming convention: `adaptSdk<SDK_MODEL_NAME>ToDatasource<DATASOURCE_MODEL_NAME>`.
So the function to adapt an SDK Instance to an Instance Datasource would be
named `adaptSdkInstanceToDatasourceInstance`.

## Resources

### Resource files

All code pertaining to resources, including models, belongs in the resource
file. The format for this file is `resource_<ENDPOINT>.go`. Ie: for
instances the filename would be `resource_instance.go`

### Resource structs

Structs should adhere to the following naming convention: `resource<NAME>`.
So for `instance` this would be `resourceInstance`.

### Resource models

Datasource model structs should adhere to the following convention:

`resourceModel<MODEL_NAME>`. So the `instance` data model would be named
`resourceModelInstance`

#### Resource adaptation functions

Adapt functions to convert SDK models to resource models have the following
naming convention: `adaptSdk<SDK_MODEL_NAME>ToResource<Resource_MODEL_NAME>`.
So the function to adapt an SDK Instance to an Instance Resource would be
named `adaptSdkInstanceToResourceInstance`.

## Validators

As validators are often shared between resource they belong in the `validators.go`
file.
