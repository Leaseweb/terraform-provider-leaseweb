package leaseweb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	LSW "github.com/LeaseWeb/leaseweb-go-sdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	leasewebAPIURL   string
	leasewebAPIToken string
	leasewebClient   *http.Client
)

// Payload -
type Payload map[string]interface{}

// Job -
type Job struct {
	UUID    string
	Status  string
	Payload Payload
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

func getLatestInstallationJob(ctx context.Context, serverID string) (*Job, error) {
	apiCtx := fmt.Sprintf("getting latest installation job for server %s", serverID)

	u, err := url.Parse(fmt.Sprintf("%s/bareMetals/v2/servers/%s/jobs", leasewebAPIURL, serverID))
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("type", "install")
	u.RawQuery = v.Encode()

	url := u.String()
	method := http.MethodGet

	response, err := doAPIRequest(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := parseErrorInfo(response.Body, apiCtx)
		logAPIError(ctx, method, url, err)
		return nil, err
	}

	var jobs struct {
		Jobs []Job
	}

	err = json.NewDecoder(response.Body).Decode(&jobs)
	if err != nil {
		return nil, NewDecodingError(apiCtx, err)
	}

	return &jobs.Jobs[0], nil
}

func getAllServers(ctx context.Context, site string) ([]LSW.DedicatedServer, error) {
	var allServers []LSW.DedicatedServer
	offset := 0
	limit := 20

	for {
		result, err := LSW.DedicatedServerApi{}.List(offset, limit, "", site)
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
