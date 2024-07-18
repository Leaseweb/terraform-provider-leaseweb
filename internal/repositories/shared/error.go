package shared

import (
	"fmt"
	"net/http"
)

type RepositoryError struct {
	msg             string
	err             error
	SdkHttpResponse *http.Response
}

func (e RepositoryError) Error() string {
	return e.msg
}

func NewSdkError(
	errorPrefix string,
	sdkError error,
	sdkHttpResponse *http.Response,
) *RepositoryError {
	return &RepositoryError{
		msg:             fmt.Errorf("%s: %w", errorPrefix, sdkError).Error(),
		err:             sdkError,
		SdkHttpResponse: sdkHttpResponse,
	}
}

func NewGeneralError(errorPrefix string, err error) *RepositoryError {
	return &RepositoryError{
		msg: fmt.Errorf("%s: %w", errorPrefix, err).Error(),
		err: err,
	}
}
