package shared

import (
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

func Test_pagination_canIncrement(t *testing.T) {
	t.Run(
		"can not increment when Offset is equal to TotalCount",
		func(t *testing.T) {
			metadata := publicCloud.Metadata{
				TotalCount: 7,
				Offset:     7,
				Limit:      2,
			}

			pagination := NewPagination(metadata)

			assert.False(t, pagination.CanIncrement())
		},
	)

	t.Run(
		"can not increment when Offset + limit passes TotalCount",
		func(t *testing.T) {
			metadata := publicCloud.Metadata{
				TotalCount: 7,
				Offset:     3,
				Limit:      10,
			}

			pagination := NewPagination(metadata)

			assert.False(t, pagination.CanIncrement())
		},
	)

	t.Run(
		"can increment when Offset + limit is less than TotalCount",
		func(t *testing.T) {
			metadata := publicCloud.Metadata{
				TotalCount: 10,
				Offset:     7,
				Limit:      2,
			}

			pagination := NewPagination(metadata)

			assert.True(t, pagination.CanIncrement())
		},
	)
}

func Test_pagination_nextPage(t *testing.T) {
	t.Run(
		"calling NextPage returns an error if we can't increment",
		func(t *testing.T) {
			metadata := publicCloud.Metadata{
				TotalCount: 10,
				Offset:     11,
				Limit:      2,
			}

			pagination := NewPagination(metadata)

			err := pagination.NextPage()

			assert.Error(t, err)
			assert.ErrorContains(
				t,
				err,
				"cannot increment as next offset 21 is larger than the total 10",
			)

		},
	)

	t.Run(
		"NextPage increments successfully",
		func(t *testing.T) {
			metadata := publicCloud.Metadata{
				TotalCount: 10,
				Offset:     5,
				Limit:      2,
			}

			pagination := NewPagination(metadata)

			err := pagination.NextPage()

			assert.NoError(t, err)
			assert.Equal(t, 7, pagination.Offset)
		},
	)
}
