package leaseweb

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	LSW "github.com/LeaseWeb/leaseweb-go-sdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	leasewebAPIURL   string
	leasewebAPIToken string
	leasewebClient   *http.Client
)

// Server -
type Server struct {
	ID string
}

// ErrorInfo -
type ErrorInfo struct {
	Context       string
	CorrelationID string              `json:"correlationId"`
	Code          string              `json:"errorCode"`
	Message       string              `json:"errorMessage"`
	Details       map[string][]string `json:"errorDetails"`
}

func (erri *ErrorInfo) Error() string {
	return "(" + erri.Code + ") " + erri.Context + ": " + erri.Message
}

// DecodingError -
type DecodingError struct {
	Context string
	Message string
}

func (errd *DecodingError) Error() string {
	return errd.Context + ": error while decoding JSON response body (" + errd.Message + ")"
}

// NewDecodingError -
func NewDecodingError(ctx string, err error) *DecodingError {
	return &DecodingError{Context: ctx, Message: err.Error()}
}

// EncodingError -
type EncodingError struct {
	Context string
	Message string
}

func (erre *EncodingError) Error() string {
	return erre.Context + ": error while encoding JSON request body (" + erre.Message + ")"
}

// NewEncodingError -
func NewEncodingError(ctx string, err error) *EncodingError {
	return &EncodingError{Context: ctx, Message: err.Error()}
}

func parseErrorInfo(r io.Reader, ctx string) error {
	erri := ErrorInfo{Context: ctx}

	if err := json.NewDecoder(r).Decode(&erri); err != nil {
		return NewDecodingError(ctx, err)
	}

	return &erri
}

func doAPIRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Lsw-Auth", leasewebAPIToken)

	if method == http.MethodPost || method == http.MethodPut {
		// not always needed even for those methods but this is simpler for now
		request.Header.Set("Content-Type", "application/json")
	}

	tflog.Trace(ctx, "executing API request", map[string]interface{}{
		"url":    url,
		"method": method,
	})

	response, err := leasewebClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func logAPIError(ctx context.Context, method, url string, err error) {
	fields := map[string]interface{}{
		"url":    url,
		"method": method,
	}

	if erri, ok := err.(*ErrorInfo); ok {
		fields["context"] = erri.Context
		fields["code"] = erri.Code
		fields["message"] = erri.Message
		fields["correlation_id"] = erri.CorrelationID

		if len(erri.Details) != 0 {
			for field, details := range erri.Details {
				fields["detail_"+field] = details
			}
		}
	} else {
		fields["message"] = err.Error()
	}

	tflog.Error(ctx, "API request error", fields)
}

func LogApiError(ctx context.Context, err error) {
	fields := map[string]interface{}{}

	if erra, ok := err.(*LSW.ApiError); ok {
		fields["url"] = erra.Url
		fields["method"] = erra.Method
		fields["code"] = erra.Code
		fields["message"] = erra.Message
		fields["correlation_id"] = erra.CorrelationId

		if len(erra.Details) != 0 {
			for field, details := range erra.Details {
				fields["detail_"+field] = details
			}
		}
	} else {
		fields["message"] = err.Error()

		if errd, ok := err.(*LSW.DecodingError); ok {
			fields["url"] = errd.Url
			fields["method"] = errd.Method
		} else if erre, ok := err.(*LSW.EncodingError); ok {
			fields["url"] = erre.Url
			fields["method"] = erre.Method
		}
	}

	tflog.Error(ctx, "API request error", fields)
}

func getAllServers(ctx context.Context, site string) ([]LSW.DedicatedServer, error) {
	var allServers []LSW.DedicatedServer
	offset := 0
	limit := 20

	for {

		opts := LSW.DedicatedServerListOptions{
			PaginationOptions: LSW.PaginationOptions{
				Offset: &offset,
				Limit:  &limit,
			},
			Site: &site,
		}

		result, err := LSW.DedicatedServerApi{}.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		if len(result.Servers) == 0 {
			break
		}

		allServers = append(allServers, result.Servers...)
		offset += limit
	}

	return allServers, nil
}
