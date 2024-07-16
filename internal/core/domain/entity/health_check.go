package entity

import (
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

type HealthCheck struct {
	Method enum.Method
	Uri    string
	Host   *string
	Port   int
}

type OptionalHealthCheckValues struct {
	Host *string
}

func NewHealthCheck(
	method enum.Method,
	uri string,
	port int,
	options OptionalHealthCheckValues,
) HealthCheck {
	healthCheck := HealthCheck{Method: method, Uri: uri, Port: port}

	healthCheck.Host = options.Host

	return healthCheck
}
