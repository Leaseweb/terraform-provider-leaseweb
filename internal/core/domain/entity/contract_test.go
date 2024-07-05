package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"terraform-provider-leaseweb/internal/core/shared/value_object/enum"
)

func TestNewContract(t *testing.T) {
	t.Run("required values are set", func(t *testing.T) {
		renewalsAt := time.Now()
		createdAt := time.Now()
		endsAt := time.Now()

		got, err := NewContract(
			enum.ContractBillingFrequencySix,
			enum.ContractTermThree,
			enum.ContractTypeMonthly,
			renewalsAt,
			createdAt,
			enum.Active,
			&endsAt,
		)

		assert.Nil(t, err)
		assert.Equal(
			t,
			enum.ContractBillingFrequencySix,
			got.BillingFrequency,
		)
		assert.Equal(t, enum.ContractTermThree, got.Term)
		assert.Equal(t, enum.ContractTypeMonthly, got.Type)
		assert.Equal(t, renewalsAt, got.RenewalsAt)
		assert.Equal(t, createdAt, got.CreatedAt)
		assert.Equal(t, enum.Active, got.State)
		assert.Equal(t, endsAt, *got.EndsAt)
	})

	t.Run(
		"error is returned when contract type is monthly and contract term is zero",
		func(t *testing.T) {
			_, err := NewContract(
				enum.ContractBillingFrequencySix,
				enum.ContractTermZero,
				enum.ContractTypeMonthly,
				time.Now(),
				time.Now(),
				enum.Active,
				nil,
			)

			assert.NotNil(t, err)
		},
	)

	t.Run(
		"error is returned when contract type is hourly and contract term is not zero",
		func(t *testing.T) {
			_, err := NewContract(
				enum.ContractBillingFrequencySix,
				enum.ContractTermThree,
				enum.ContractTypeHourly,
				time.Now(),
				time.Now(),
				enum.Active,
				nil,
			)

			assert.NotNil(t, err)
		},
	)

}
