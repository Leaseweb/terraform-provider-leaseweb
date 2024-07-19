package shared

import (
	"fmt"

	"terraform-provider-leaseweb/internal/repositories/shared"
	shared2 "terraform-provider-leaseweb/internal/shared"
)

type ServiceError struct {
	msg           string
	GeneralError  error
	ErrorResponse *shared2.ErrorResponse
}

func (e ServiceError) Error() string {
	return e.msg
}

func NewFromRepositoryError(
	errorPrefix string,
	repositoryError *shared.RepositoryError,
) *ServiceError {
	return &ServiceError{
		msg:           fmt.Errorf("%s: %w", errorPrefix, repositoryError).Error(),
		ErrorResponse: repositoryError.ErrorResponse,
	}
}

func NewError(
	errorPrefix string,
	err error,
) *ServiceError {
	return &ServiceError{
		msg:          fmt.Errorf("%s: %w", errorPrefix, err).Error(),
		GeneralError: err,
	}
}
