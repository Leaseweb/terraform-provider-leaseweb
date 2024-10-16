// Package errors implements errors related to core services.
package errors

import (
	"fmt"

	repository2 "github.com/leaseweb/terraform-provider-leaseweb/internal/provider/shared/repository"
)

type ServiceError struct {
	msg           string
	GeneralError  error
	ErrorResponse *repository2.ErrorResponse
}

func (e ServiceError) Error() string {
	return e.msg
}

// NewFromRepositoryError generates a new error from the passed repository error.
func NewFromRepositoryError(
	errorPrefix string,
	repositoryError repository2.RepositoryError,
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