package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Error struct {
	err  error
	resp *http.Response
}

// TODO: we need to merge error.go and log_error.go and have unified error/logging functionality.
func (e Error) Error() string {
	// Check if the response or its body is nil,
	//or if the status code is not an error.
	if e.resp == nil || e.resp.Body == nil || e.resp.StatusCode < 400 {
		return e.err.Error()
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}(e.resp.Body)

	// Try to decode the response body as JSON.
	var errorResponse map[string]interface{}
	if err := json.NewDecoder(e.resp.Body).Decode(&errorResponse); err != nil {
		return e.err.Error()
	}

	// Build the error message from errorResponse.
	if msg := buildErrorMessage(errorResponse); msg != "" {
		return msg
	}

	// Default to the original error message if no relevant information is found.
	return e.err.Error()
}

// Helper function to build an error message from the decoded JSON.
func buildErrorMessage(errorResponse map[string]interface{}) string {
	var msg string

	// Append the main error message if available.
	if errorMessage, ok := errorResponse["errorMessage"]; ok {
		msg += fmt.Sprintf("%v", errorMessage)
	}

	// Append details if available.
	if errorDetails, ok := errorResponse["errorDetails"].(map[string]interface{}); ok {
		msg += "\n" + formatErrorDetails(errorDetails)
	}

	return msg
}

// Helper function to format error details.
func formatErrorDetails(errorDetails map[string]interface{}) string {
	var detailsMsg string

	for key, details := range errorDetails {
		detailsMsg += fmt.Sprintf("%s:\n", key)

		// Check if the details are a list of messages.
		if detailList, ok := details.([]interface{}); ok {
			for _, detail := range detailList {
				detailsMsg += fmt.Sprintf("\t%s\n", detail)
			}
		}
	}

	return detailsMsg
}

func NewError(resp *http.Response, err error) Error {
	return Error{
		resp: resp,
		err:  err,
	}
}

// normalizeErrorResponseKey converts an api key path to a string
// that SetAttributeErrorsFromServerResponse can handle.
// `instanceId` & `instance.id` both become `instance_id`.
func normalizeErrorResponseKey(key string) string {
	// Assume that the key has the format `contract.id`
	//if any dots are found.
	if strings.Contains(key, ".") {
		return strings.ToLower(strings.Replace(key, ".", "_", -1))
	}

	// If no dots are found, assume camel case.
	m := regexp.MustCompile("[A-Z]")
	res := m.ReplaceAllStringFunc(key, func(s string) string {
		return "_" + s
	})

	return strings.ToLower(res)
}

// SetAttributeErrorsFromServerResponse takes a server response and maps errors to the appropriate attributes.
// If an attribute cannot be found,
// the error is shown to the user on a resource level.
func SetAttributeErrorsFromServerResponse(
	summary string,
	response *http.Response,
	diags *diag.Diagnostics,
) {
	// Nothing to do when response does not exist.
	if response == nil {
		return
	}

	// Convert server response to ErrorResponse object.
	errorResponse, err := newErrorResponse(response.Body)
	// If error cannot be translated,
	// Terraform will show a general error to the user.
	if err != nil {
		return
	}

	// Convert key returned from api to an attribute path.
	// I.e.: []string{"image", "id"}.
	for errorKey, errorDetailList := range errorResponse.ErrorDetails {
		normalizedErrorKey := normalizeErrorResponseKey(errorKey)
		mapKeys := strings.Split(normalizedErrorKey, "_")
		attributePath := path.Root(mapKeys[0])

		// Every element in the map goes one level deeper.
		for _, mapKey := range mapKeys[1:] {
			attributePath = attributePath.AtMapKey(mapKey)
		}

		// Each attribute can have multiple errors.
		for _, errorDetail := range errorDetailList {
			diags.AddAttributeError(attributePath, summary, errorDetail)
		}
	}
}

type ErrorResponse struct {
	CorrelationId string              `json:"correlationId"`
	ErrorCode     string              `json:"errorCode"`
	ErrorMessage  string              `json:"errorMessage"`
	ErrorDetails  map[string][]string `json:"errorDetails"`
}

// newErrorResponse generates a new ErrorResponse object from an api response body.
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

// LogError tries to parse the errors from the httpResponse.
// If possible,
// it prints an error log with the error response as an additional field.
// Otherwise, it prints an error log without the response.
func LogError(ctx context.Context, httpResponse *http.Response, summary string) {
	errorResponse, err := newErrorResponse(httpResponse.Body)
	if err == nil {
		tflog.Error(ctx, summary)
	}

	tflog.Error(
		ctx,
		summary,
		map[string]interface{}{"ErrorResponse": errorResponse},
	)
}
