package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getHttpErrorMessage(resp *http.Response, err error) string {
	if resp == nil || resp.Body == nil || resp.StatusCode < 400 {
		return err.Error()
	}

	defer resp.Body.Close() // Ensure the body is closed
	var errorResponse map[string]interface{}

	// Attempt to decode the response body as JSON
	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
		if errorMessage, ok := errorResponse["errorMessage"]; ok {
			return fmt.Sprintf("%v", errorMessage)
		}
	}

	return err.Error()
}
