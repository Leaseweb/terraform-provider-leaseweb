package utils

import (
	"encoding/json"
	"net/http"
)

const DefaultErrMsg = "An error has occurred in the program. Please consider opening an issue."

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

	if msg := buildErrorMessage(errorResponse); msg != "" && msg != "{}" {
		return msg
	}

	return e.err.Error()
}

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
