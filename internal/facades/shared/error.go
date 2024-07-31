package shared

import (
	"fmt"

	sharedService "terraform-provider-leaseweb/internal/core/services/errors"
	"terraform-provider-leaseweb/internal/repositories/shared"
)

type FacadeError struct {
	msg           string
	ErrorResponse *shared.ErrorResponse
	GeneralError  error
}

func (e FacadeError) Error() string {
	return e.msg
}

// NewFromServicesError generates a new facade error from a ServiceError.
func NewFromServicesError(
	errorPrefix string,
	serviceError *sharedService.ServiceError,
) *FacadeError {
	return &FacadeError{
		msg:           fmt.Errorf("%s: %w", errorPrefix, serviceError).Error(),
		GeneralError:  serviceError.GeneralError,
		ErrorResponse: serviceError.ErrorResponse,
	}
}

// NewError generates a regular facade error.
func NewError(
	errorPrefix string,
	err error,
) *FacadeError {
	return &FacadeError{
		msg:          fmt.Errorf("%s: %w", errorPrefix, err).Error(),
		GeneralError: err,
	}
}
