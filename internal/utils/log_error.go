package utils

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
func newErrorResponse(body io.Reader) (*ErrorResponse, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, body)
	if err != nil {
		return nil, err
	}

	errorResponse := ErrorResponse{}

	jsonErr := json.Unmarshal([]byte(buf.String()), &errorResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &errorResponse, nil
}

type SdkError struct {
	msg           string
	err           error
	ErrorResponse *ErrorResponse
}

func (s SdkError) Error() string {
	return s.msg
}

// NewSdkError generates a new error from an sdk error & response.
func NewSdkError(
	errorPrefix string,
	sdkError error,
	sdkHttpResponse *http.Response,
) *SdkError {
	repositoryError := SdkError{
		msg: fmt.Errorf("%s: %w", errorPrefix, sdkError).Error(),
		err: sdkError,
	}

	// Convert the returned JSON to an ErrorResponse struct.
	if sdkHttpResponse != nil {
		errorResponse, err := newErrorResponse(sdkHttpResponse.Body)
		if err == nil {
			repositoryError.ErrorResponse = errorResponse
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
