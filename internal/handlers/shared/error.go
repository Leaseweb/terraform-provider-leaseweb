package shared

import (
	"fmt"
	"io"
	"strings"

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

func (e HandlerError) GetResponse() *string {
	if e.ServiceError != nil {
		if e.ServiceError.RepositoryError != nil {
			if e.ServiceError.RepositoryError.SdkHttpResponse != nil {
				body := e.ServiceError.RepositoryError.SdkHttpResponse.Body
				buf := new(strings.Builder)
				_, sdkResponseError := io.Copy(buf, body)
				if sdkResponseError == nil {
					bodyContent := buf.String()
					return &bodyContent
				}
			}
		}
	}

	return nil
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
