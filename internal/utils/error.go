package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	err  error
	resp *http.Response
}

func (e Error) Error() string {
	if e.resp == nil || e.resp.Body == nil || e.resp.StatusCode < 400 {
		return e.err.Error()
	}

	defer e.resp.Body.Close() // Ensure the body is closed
	var errorResponse map[string]interface{}

	// Attempt to decode the response body as JSON
	if err := json.NewDecoder(e.resp.Body).Decode(&errorResponse); err == nil {
		if errorMessage, ok := errorResponse["errorMessage"]; ok {
			return fmt.Sprintf("%v", errorMessage)
		}
	}

	return e.err.Error()
}

func NewError(resp *http.Response, err error) Error {
	return Error{
		resp: resp,
		err:  err,
	}
}
