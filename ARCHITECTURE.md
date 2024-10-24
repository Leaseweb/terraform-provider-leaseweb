# Code architecture

## Packages

Each package corresponds to a product group, i.e.: `publiccloud`

## Datasources

All code pertaining to datasources, including models, belongs in the datasource
file. The format for this file is `<PRODUCT>_<ENDPOINT>_datasource.go`. Ie: for
instances the filename would be `publiccloud_instances_datasource.go`

### Datasource models

Datasource model structs should adhere to the following convention:

`datasourceModel<ModelName>`. So the instances data model would be called
`dataSourceModelinstances`

## Resources

All code pertaining to resources, including models, belongs in the resource
file. The format for this file is `<PRODUCT>_<ENDPOINT>_resource.go`. Ie: for
instances the filename would be `publiccloud_instances_resource.go`

### Resource models

Datasource model structs should adhere to the following convention:

`resourceModel<ModelName>`. So the instances data model would be called
`resourceModelinstance`

## Validators

As validators are often shared between resource they belong in the `validators.go`
file.
