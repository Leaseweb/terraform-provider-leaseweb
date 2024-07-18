package public_cloud_repository

import (
	"fmt"
	"net/http"
)

type PublicCloudRepositoryError struct {
	msg             string
	err             error
	SdkHttpResponse *http.Response
}

func (e PublicCloudRepositoryError) Error() string {
	return e.msg
}

func newPublicCloudRepositoryError(
	errorPrefix string,
	sdkError error,
	sdkHttpResponse *http.Response,
) error {
	return PublicCloudRepositoryError{
		msg:             fmt.Errorf("%s: %w", errorPrefix, sdkError).Error(),
		err:             sdkError,
		SdkHttpResponse: sdkHttpResponse,
	}
}
