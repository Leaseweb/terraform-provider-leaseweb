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

func NewRepositoryError(
	errorPrefix string,
	sdkError error,
	sdkHttpResponse *http.Response,
) error {
	return RepositoryError{
		msg:             fmt.Errorf("%s: %w", errorPrefix, sdkError).Error(),
		err:             sdkError,
		SdkHttpResponse: sdkHttpResponse,
	}
}
