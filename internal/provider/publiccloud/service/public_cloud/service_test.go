package public_cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_ValidateContractTerm(t *testing.T) {
	t.Run(
		"ErrContractTermCannotBeZero is returned when contract term is monthly and contract term is 0",
		func(t *testing.T) {
			service := Service{}
			got := service.ValidateContractTerm(0, "MONTHLY")

			assert.ErrorIs(t, got, ErrContractTermCannotBeZero)
		},
	)

	t.Run(
		"ErrContractTermMustBeZero is returned when contract term is hourly and contract term is not 0",
		func(t *testing.T) {
			service := Service{}
			got := service.ValidateContractTerm(3, "HOURLY")

			assert.ErrorIs(t, got, ErrContractTermMustBeZero)
		},
	)

	t.Run("no error is returned when contract is valid", func(t *testing.T) {
		service := Service{}
		got := service.ValidateContractTerm(0, "HOURLY")

		assert.Nil(t, got)
	},
	)

	t.Run(
		"error is returned when invalid contractTerm is passed",
		func(t *testing.T) {
			service := Service{}
			got := service.ValidateContractTerm(55, "HOURLY")

			assert.ErrorContains(t, got, "55")
		},
	)

	t.Run(
		"error is returned when invalid contractType is passed",
		func(t *testing.T) {
			service := Service{}
			got := service.ValidateContractTerm(0, "tralala")

			assert.ErrorContains(t, got, "tralala")
		},
	)
}
