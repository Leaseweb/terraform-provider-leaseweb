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

const DefaultErrMsg = "An error has occurred in the program. Please consider opening an issue."

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

	if msg := buildErrorMessage(errorResponse); msg != "" && msg != "{}" {
		return msg
	}

	// Default to the original error message if no relevant information is found.
	return e.err.Error()
}

// Helper function to build an error message from the decoded JSON.
func buildErrorMessage(errorResponse map[string]interface{}) string {
	// Create a map to store only the fields we care about.
	output := make(map[string]interface{})

	if errorCode, ok := errorResponse["errorCode"]; ok {
		output["errorCode"] = errorCode
	}
	if errorMessage, ok := errorResponse["errorMessage"]; ok {
		output["errorMessage"] = errorMessage
	}
	if userMessage, ok := errorResponse["userMessage"]; ok {
		output["userMessage"] = userMessage
	}
	if correlationId, ok := errorResponse["correlationId"]; ok {
		output["correlationId"] = correlationId
	}
	if errorDetails, ok := errorResponse["errorDetails"]; ok {
		output["errorDetails"] = errorDetails
	}

	// Encode the output map as a JSON string.
	jsonOutput, _ := json.MarshalIndent(output, "", "  ")
	return string(jsonOutput)
}

func NewError(resp *http.Response, err error) Error {
	return Error{
		resp: resp,
		err:  err,
	}
}

// normalizeErrorResponseKey converts an api key path to a string
// that HandleSdkError can handle.
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

// HandleSdkError takes a server response & error
// and maps errors to the appropriate attributes.
// If an attribute cannot be found,
// the error is shown to the user on a resource level.
// A DEBUG log is also created with all the relevant information.
func HandleSdkError(
	summary string,
	httpResponse *http.Response,
	err error,
	diags *diag.Diagnostics,
	ctx context.Context,
) {
	// Nothing to do when httpResponse does not exist.
	if httpResponse == nil {
		handleError(summary, err, diags)
		return
	}

	// Try to read httpResponse body into buffer
	buf := new(strings.Builder)
	_, err = io.Copy(buf, httpResponse.Body)
	if err != nil {
		handleError(summary, err, diags)
		return
	}

	// Create DEBUG log with httpResponse body.
	responseMap, err := newResponseMap(buf.String())
	if err != nil {
		tflog.Debug(
			ctx,
			summary,
			map[string]any{"httpResponse": fmt.Sprintf("%v", httpResponse.Body)},
		)
		handleError(summary, nil, diags)
		return
	}
	tflog.Debug(ctx, summary, map[string]any{"response": responseMap})

	// Convert httpResponse buffer to ErrorResponse object.
	errorResponse, err := newErrorResponse(buf.String())
	if err != nil {
		handleError(summary, err, diags)
		return
	}

	// Convert key returned from api to an attribute path.
	// I.e.: []string{"image", "id"}.
	errorSet := false
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
		errorSet = true
	}

	// If no attribute errors are set, set a global error
	if !errorSet {
		errorResponseString, err := json.MarshalIndent(
			errorResponse,
			"",
			" ",
		)
		if err != nil {
			handleError(summary, err, diags)
			return
		}

		diags.AddError(summary, string(errorResponseString))
	}
}

type ErrorResponse struct {
	CorrelationId string              `json:"correlationId,omitempty"`
	ErrorCode     string              `json:"errorCode,omitempty"`
	ErrorMessage  string              `json:"errorMessage,omitempty"`
	ErrorDetails  map[string][]string `json:"errorDetails,omitempty"`
}

func handleError(
	summary string,
	err error,
	diags *diag.Diagnostics,
) {
	if err != nil {
		diags.AddError(summary, err.Error())
		return
	}

	diags.AddError(summary, DefaultErrMsg)
}

func newErrorResponse(body string) (*ErrorResponse, error) {
	errorResponse := ErrorResponse{}

	jsonErr := json.Unmarshal([]byte(body), &errorResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &errorResponse, nil
}

func newResponseMap(body string) (map[string]interface{}, error) {
	var response map[string]interface{}

	jsonErr := json.Unmarshal([]byte(body), &response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response, nil
}
