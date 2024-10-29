package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

type Error struct {
	err  error
	resp *http.Response
}

// TODO: we need to merge error.go and log_error.go and have unified error/logging functionality.
func (e Error) Error() string {
	// Check if response or its body is nil, or if the status code is not an error.
	if e.resp == nil || e.resp.Body == nil || e.resp.StatusCode < 400 {
		return e.err.Error()
	}
	defer e.resp.Body.Close()

	// Try to decode the response body as JSON.
	var errorResponse map[string]interface{}
	if err := json.NewDecoder(e.resp.Body).Decode(&errorResponse); err != nil {
		return e.err.Error()
	}

	// Build the error message from errorResponse.
	if msg := buildErrorMessage(errorResponse); msg != "" {
		return msg
	}

	// Default to original error message if no relevant information is found.
	return e.err.Error()
}

// Helper function to build error message from the decoded JSON.
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

// SetAttributeErrorsFromServerResponse takes a server response and maps errors to the appropriate attributes.
// If an attribute cannot be found,
// the error is shown to the user on a resource level.
func SetAttributeErrorsFromServerResponse(
	summary string,
	response http.Response,
	diags *diag.Diagnostics,
) {
	// Convert server response to ErrorResponse object.
	errorResponse, err := newErrorResponse(response.Body)
	// If error cannot be translated,
	// Terraform will show a general error to the user.
	if err != nil {
		return
	}

	// Convert key
	// returned from api to an attribute path
	// by splitting up the cameCase to a map of lowercase words
	// []string{"camel", "case"}.
	for errorKey, errorDetailList := range errorResponse.ErrorDetails {
		m := regexp.MustCompile("[A-Z]")
		res := m.ReplaceAllStringFunc(errorKey, func(s string) string {
			return "_" + s
		})

		res = strings.ToLower(res)

		mapKeys := strings.Split(res, "_")
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
