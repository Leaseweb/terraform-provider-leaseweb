package shared

import (
	"fmt"

	sharedService "terraform-provider-leaseweb/internal/core/services/shared"
	"terraform-provider-leaseweb/internal/shared"
)

type HandlerError struct {
	msg           string
	ErrorResponse *shared.ErrorResponse
	GeneralError  error
}

func (e HandlerError) Error() string {
	return e.msg
}

// NewFromServicesError generates a new handler error from a ServiceError.
func NewFromServicesError(
	errorPrefix string,
	serviceError *sharedService.ServiceError,
) *HandlerError {
	return &HandlerError{
		msg:           fmt.Errorf("%s: %w", errorPrefix, serviceError).Error(),
		GeneralError:  serviceError.GeneralError,
		ErrorResponse: serviceError.ErrorResponse,
	}
}

// NewError generates a regular handler error.
func NewError(
	errorPrefix string,
	err error,
) *HandlerError {
	return &HandlerError{
		msg:          fmt.Errorf("%s: %w", errorPrefix, err).Error(),
		GeneralError: err,
	}
}
