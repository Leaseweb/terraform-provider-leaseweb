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

const defaultErrMsg = "An error has occurred in the program. Please consider opening an issue."

// GeneralError should be called when general errors need to be handled.
func GeneralError(diags *diag.Diagnostics, ctx context.Context, err error) {
	if err != nil {
		tflog.Debug(ctx, err.Error())
	}
	diags.AddError("Unexpected Error", defaultErrMsg)
}

// ImportOnlyError should be used in resource Read() functions for resources that can only be imported.
func ImportOnlyError(diags *diag.Diagnostics) {
	diags.AddError("Resource can only be imported, not created.", "")
}

// ConfigError should be used in resource/data source Configure() functions to handle unexpected resource configure types.
func ConfigError(diags *diag.Diagnostics, providerData any) {
	diags.AddError(
		"Unexpected Resource Configure Type",
		fmt.Sprintf(
			"Expected client.Client, got: %T. Please report this issue to the provider developers.",
			providerData,
		),
	)
}

// UnexpectedImportIdentifierError should be used in Import() functions where the identifier is incorrect.
func UnexpectedImportIdentifierError(diags *diag.Diagnostics, format string, got string) {
	diags.AddError(
		"Unexpected Import Identifier",
		fmt.Sprintf(
			"Expected import identifier with format: %q. Got: %q",
			format,
			got,
		),
	)
}

// SdkError should be used to handle errors returned by the SDK.
func SdkError(
	ctx context.Context,
	diags *diag.Diagnostics,
	err error,
	resp *http.Response,
) {
	errorHandler := errorHandler{
		ctx:     ctx,
		diags:   diags,
		summary: "Unexpected API Error",
		err:     err,
		resp:    resp,
	}

	errorHandler.report()
}

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

func (e *errorHandler) report() {
	if e.diags == nil {
		e.writeLog("unable to record error details.")
		log.Fatal(e.summary, defaultErrMsg)
	}

	if e.err == nil {
		e.writeLog("No error detail found.")
		e.writeOutput(defaultErrMsg)
		return
	}

	if e.resp != nil && e.resp.Body != nil && e.resp.StatusCode >= 400 {
		e.handleHTTPError()
		return
	}

	e.writeOutput(e.err.Error())
}

func (e *errorHandler) handleHTTPError() {
	if e.resp.StatusCode == 504 {
		e.writeLog(fmt.Sprintf("server response: %v", e.resp.Body))
		e.writeOutput("The server took too long to respond.")
		return
	}

	if e.resp.StatusCode == 404 {
		e.writeLog(fmt.Sprintf("server response: %v", e.resp.Body))
		e.writeOutput("Resource not found.")
		return
	}

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
		e.writeOutput(defaultErrMsg)
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
		e.writeLog(fmt.Sprintf(
			"failed to format error response as JSON: %v",
			err,
		))
		e.writeOutput(defaultErrMsg)
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
