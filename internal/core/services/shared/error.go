package shared

import (
	"fmt"

	repository "terraform-provider-leaseweb/internal/repositories/shared"
	"terraform-provider-leaseweb/internal/shared"
)

type ServiceError struct {
	msg           string
	GeneralError  error
	ErrorResponse *shared.ErrorResponse
}

func (e ServiceError) Error() string {
	return e.msg
}

// NewFromRepositoryError generates a new error from the passed repository error.
func NewFromRepositoryError(
	errorPrefix string,
	repositoryError repository.RepositoryError,
) *ServiceError {
	return &ServiceError{
		msg:           fmt.Errorf("%s: %w", errorPrefix, repositoryError).Error(),
		ErrorResponse: repositoryError.ErrorResponse,
	}
}

// NewError generates a new general error.
func NewError(
	errorPrefix string,
	err error,
) *ServiceError {
	return &ServiceError{
		msg:          fmt.Errorf("%s: %w", errorPrefix, err).Error(),
		GeneralError: err,
	}
}
