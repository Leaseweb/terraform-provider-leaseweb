package utils

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

		want := "{\n  \"errorMessage\": \"this is error message\"\n}"
		got := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, want, got)
	})

	t.Run("errorMessage must return if response contains object of errorMessage", func(t *testing.T) {
		errorContent, _ := json.Marshal(map[string]interface{}{"errorMessage": map[string]string{"a": "b", "c": "d"}})
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(string(errorContent))),
		}

		want := "{\n  \"errorMessage\": {\n    \"a\": \"b\",\n    \"c\": \"d\"\n  }\n}"
		got := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, want, got)
	})

	t.Run("errorMessage must return if response with error details if they exists", func(t *testing.T) {
		errorContent, _ := json.Marshal(map[string]interface{}{
			"errorMessage": map[string]string{"a": "b", "c": "d"},
			"errorDetails": map[string][]string{
				"password": {
					"this value should not be blank",
					"blah blah",
				},
				"email": {
					"this value should be valid",
					"blah2 blah2",
				},
			},
		})
		resp := &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(string(errorContent))),
		}

		want := "{\n  \"errorDetails\": {\n    \"email\": [\n      \"this value should be valid\",\n      \"blah2 blah2\"\n    ],\n    \"password\": [\n      \"this value should not be blank\",\n      \"blah blah\"\n    ]\n  },\n  \"errorMessage\": {\n    \"a\": \"b\",\n    \"c\": \"d\"\n  }\n}"
		got := NewError(resp, errors.New("error content")).Error()
		assert.Equal(t, want, got)
	})
}
