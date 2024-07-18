package shared

import (
	"fmt"

	sharedService "terraform-provider-leaseweb/internal/core/services/shared"
)

type HandlerError struct {
	msg          string
	ServiceError *sharedService.ServiceError
	GeneralError error
}

func (e HandlerError) Error() string {
	return e.msg
}

func NewServiceError(
	errorPrefix string,
	serviceError *sharedService.ServiceError,
) *HandlerError {
	return &HandlerError{
		msg:          fmt.Errorf("%s: %w", errorPrefix, serviceError).Error(),
		ServiceError: serviceError,
	}
}

func NewGeneralError(
	errorPrefix string,
	err error,
) *HandlerError {
	return &HandlerError{
		msg:          fmt.Errorf("%s: %w", errorPrefix, err).Error(),
		GeneralError: err,
	}
}
