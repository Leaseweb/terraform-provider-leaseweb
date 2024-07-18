package shared

import (
	"fmt"

	"terraform-provider-leaseweb/internal/repositories/shared"
)

type ServiceError struct {
	msg             string
	RepositoryError *shared.RepositoryError
	GeneralError    error
}

func (e ServiceError) Error() string {
	return e.msg
}

func NewRepositoryError(
	errorPrefix string,
	repositoryError *shared.RepositoryError,
) *ServiceError {
	return &ServiceError{
		msg:             fmt.Errorf("%s: %w", errorPrefix, repositoryError).Error(),
		RepositoryError: repositoryError,
	}
}

func NewGeneralError(
	errorPrefix string,
	err error,
) *ServiceError {
	return &ServiceError{
		msg:          fmt.Errorf("%s: %w", errorPrefix, err).Error(),
		GeneralError: err,
	}
}
