package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const DefaultErrMsg = "An error has occurred in the program. Please consider opening an issue."

type errorResponse struct {
	ErrorCode     string              `json:"errorCode,omitempty"`
	ErrorMessage  string              `json:"errorMessage,omitempty"`
	UserMessage   string              `json:"userMessage,omitempty"`
	CorrelationId string              `json:"correlationId,omitempty"`
	ErrorDetails  map[string][]string `json:"errorDetails,omitempty"`
}

type errorHandler struct {
	summary string
	resp    *http.Response
	err     error
	diags   *diag.Diagnostics
	ctx     context.Context
}

func Error(
	ctx context.Context,
	diags *diag.Diagnostics,
	summary string,
	err error,
	resp *http.Response,
) {

	errorHandler := errorHandler{
		ctx:     ctx,
		diags:   diags,
		summary: summary,
		err:     err,
		resp:    resp,
	}

	errorHandler.report()
}

func (e *errorHandler) report() {

	if e.diags == nil {
		e.writeLog("unable to record error details.")
		log.Fatal(e.summary, DefaultErrMsg)
	}

	if e.err == nil {
		e.writeLog("No error detail found.")
		e.writeOutput(DefaultErrMsg)
		return
	}

	if e.resp != nil && e.resp.Body != nil && e.resp.StatusCode >= 400 {
		e.handleHTTPError()
		return
	}

	e.writeOutput(e.err.Error())
}

func (e *errorHandler) handleHTTPError() {
	// Close response body with direct defer reference for clarity
	defer func() {
		if err := e.resp.Body.Close(); err != nil {
			e.writeLog(fmt.Sprintf("error closing response body: %v", err))
		}
	}()

	e.writeLog(fmt.Sprintf("response body: %v", e.resp.Body))
	var errorResponse errorResponse
	if err := json.NewDecoder(e.resp.Body).Decode(&errorResponse); err != nil {
		e.writeLog(fmt.Sprintf("error decoding HTTP response body: %v", err))
		e.writeOutput(DefaultErrMsg)
		return
	}

	e.processErrorResponse(errorResponse)
}

// processErrorResponse checks for validation errors in ErrorDetails and handles them if present.
func (e *errorHandler) processErrorResponse(errorResponse errorResponse) {
	if len(errorResponse.ErrorDetails) > 0 {
		e.handleValidationErrors(errorResponse)
		if e.diags.HasError() {
			return
		}
	}

	// Attempt to convert errorResponse to JSON format for output
	jsonOutput, err := json.MarshalIndent(errorResponse, "", "  ")
	if err != nil {
		e.writeLog(fmt.Sprintf("failed to format error response as JSON: %v", err))
		e.writeOutput(DefaultErrMsg)
		return
	}

	e.writeOutput(string(jsonOutput))
}

func (e *errorHandler) handleValidationErrors(errorResponse errorResponse) {
	// Convert key returned from api to an attribute path.
	// I.e.: []string{"image", "id"}.
	for errorKey, errorDetailList := range errorResponse.ErrorDetails {
		normalizedErrorKey := mapErrorDetailsKey(errorKey)
		mapKeys := strings.Split(normalizedErrorKey, "_")
		attributePath := path.Root(mapKeys[0])

		// Every element in the map goes one level deeper.
		for _, mapKey := range mapKeys[1:] {
			attributePath = attributePath.AtMapKey(mapKey)
		}

		// Each attribute can have multiple errors.
		for _, errorDetail := range errorDetailList {
			e.diags.AddAttributeError(attributePath, e.summary, errorDetail)
		}
	}
}

// `instanceId` & `instance.id` both become `instance_id`.
func mapErrorDetailsKey(key string) string {
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

func (e *errorHandler) writeLog(details string) {
	tflog.Debug(e.ctx, e.summary, map[string]any{"details": details})
}

func (e *errorHandler) writeOutput(details string) {
	e.diags.AddError(e.summary, details)
}
