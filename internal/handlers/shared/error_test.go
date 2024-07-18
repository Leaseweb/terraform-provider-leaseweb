package shared

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	sharedService "terraform-provider-leaseweb/internal/core/services/shared"
	sharedRepository "terraform-provider-leaseweb/internal/repositories/shared"
)

func TestNewServiceError(t *testing.T) {
	err := errors.New("tralala")
	repositoryError := sharedRepository.NewGeneralError(
		"repositoryErrorPrefix",
		err,
	)

	serviceError := sharedService.NewRepositoryError(
		"serviceErrorPrefix",
		repositoryError,
	)
	got := NewServiceError("prefix", serviceError)

	want := HandlerError{
		msg:          "prefix: serviceErrorPrefix: repositoryErrorPrefix: tralala",
		ServiceError: serviceError,
	}

	assert.Equal(t, want, *got)
}

func TestHandlerError_Error(t *testing.T) {
	err := HandlerError{msg: "tralala"}
	want := "tralala"
	got := err.Error()

	assert.Equal(t, want, got)
}

func TestNewGeneralError(t *testing.T) {
	err := errors.New("tralala")

	got := NewGeneralError("prefix", err)

	want := HandlerError{
		msg:          "prefix: tralala",
		GeneralError: err,
	}

	assert.Equal(t, want, *got)
}

func TestHandlerError_GetResponse(t *testing.T) {
	t.Run("nil returned if serviceError is not set", func(t *testing.T) {
		handlerError := HandlerError{}
		got := handlerError.GetResponse()

		assert.Nil(t, got)
	})

	t.Run("nil returned if repositoryError is not set", func(t *testing.T) {
		handlerError := HandlerError{ServiceError: &sharedService.ServiceError{}}
		got := handlerError.GetResponse()

		assert.Nil(t, got)
	})

	t.Run("nil returned if http response is not set", func(t *testing.T) {
		handlerError := HandlerError{
			ServiceError: &sharedService.ServiceError{
				RepositoryError: &sharedRepository.RepositoryError{},
			},
		}
		got := handlerError.GetResponse()

		assert.Nil(t, got)
	})

	t.Run("response body is returned if set", func(t *testing.T) {
		handlerError := HandlerError{
			ServiceError: &sharedService.ServiceError{
				RepositoryError: &sharedRepository.RepositoryError{
					SdkHttpResponse: &http.Response{
						Body: io.NopCloser(bytes.NewReader([]byte("tralala"))),
					},
				},
			},
		}
		got := handlerError.GetResponse()

		assert.Equal(t, "tralala", *got)
	})
}
