package shared

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RepositoryError struct {
	msg           string
	err           error
	ErrorResponse *ErrorResponse
}

func (e RepositoryError) Error() string {
	return e.msg
}

// NewSdkError generates a new error from an sdk error & response.
func NewSdkError(
	errorPrefix string,
	sdkError error,
	sdkHttpResponse *http.Response,
) *RepositoryError {
	repositoryError := RepositoryError{
		msg: fmt.Errorf("%s: %w", errorPrefix, sdkError).Error(),
		err: sdkError,
	}

	// Convert the returned JSON to an ErrorResponse struct.
	if sdkHttpResponse != nil {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, sdkHttpResponse.Body)
		if err == nil {
			bodyContent := buf.String()
			errorResponse, err := NewErrorResponse(bodyContent)
			if err == nil {
				repositoryError.ErrorResponse = errorResponse
			}
		}
	}

	return &repositoryError
}
