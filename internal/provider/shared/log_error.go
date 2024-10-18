package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ErrorResponse struct {
	CorrelationId string              `json:"correlationId"`
	ErrorCode     string              `json:"errorCode"`
	ErrorMessage  string              `json:"errorMessage"`
	ErrorDetails  map[string][]string `json:"errorDetails"`
}

// newErrorResponse generates a new ErrorResponse object from an api response.
func newErrorResponse(jsonStr string) (*ErrorResponse, error) {
	errorResponse := ErrorResponse{}

	err := json.Unmarshal([]byte(jsonStr), &errorResponse)
	if err != nil {
		return nil, err
	}

	return &errorResponse, nil
}

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
			errorResponse, err := newErrorResponse(bodyContent)
			if err == nil {
				repositoryError.ErrorResponse = errorResponse
			}
		}
	}

	return &repositoryError
}

// LogError prints the passed errorResponse as a Terraform error log.
func LogError(
	ctx context.Context,
	errorResponse *ErrorResponse,
	diags *diag.Diagnostics,
	summary string,
	detail string,
) {
	if errorResponse != nil {
		tflog.Error(
			ctx,
			summary,
			map[string]interface{}{"ErrorResponse": errorResponse},
		)
	}

	diags.AddError(summary, detail)
}
