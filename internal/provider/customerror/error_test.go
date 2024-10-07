package customerror

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHttpErrorMessage(t *testing.T) {

	t.Run("error must return if response is nil", func(t *testing.T) {
		message := NewError(nil, errors.New("error content")).Error()
		assert.Equal(t, "error content", message)
	})

	t.Run("error must return if body of response is nil", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 500,
			Body:       nil,
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "error content", message)
	})

	t.Run("error must return if response is 2xx or 3xx", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("body message")),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "error content", message)
	})

	t.Run("error must return if response is not contain errorMessage", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString("body message")),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "error content", message)
	})

	t.Run("error must return if response is not contain errorMessage", func(t *testing.T) {
		errorContent, _ := json.Marshal(map[string]string{"balh": "balh"})
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(string(errorContent))),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "error content", message)
	})

	t.Run("errorMessage must return if response contains string errorMessage", func(t *testing.T) {
		errorContent, _ := json.Marshal(map[string]string{"errorMessage": "this is error message"})
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(string(errorContent))),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "this is error message", message)
	})

	t.Run("errorMessage must return if response contains object of errorMessage", func(t *testing.T) {
		errorContent, _ := json.Marshal(map[string]interface{}{"errorMessage": map[string]string{"a": "b", "c": "d"}})
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(string(errorContent))),
		}
		message := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, "map[a:b c:d]", message)
	})
}
