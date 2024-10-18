package shared

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
	"github.com/stretchr/testify/assert"
)

type exampleRequest struct {
	CurrentOffset int32
}

func (e exampleRequest) Offset(offset int32) exampleRequest {
	e.CurrentOffset = offset
	return e
}

func Test_pagination_canIncrement(t *testing.T) {
	t.Run(
		"can not increment when offset is equal to totalCount",
		func(t *testing.T) {
			request := publicCloud.ApiGetInstanceListRequest{}

			pagination := NewPagination(10, 5, request)
			pagination.offset = 5

			assert.False(t, pagination.CanIncrement())
		},
	)

	t.Run(
		"can not increment when offset + limit passes totalCount",
		func(t *testing.T) {
			request := publicCloud.ApiGetInstanceListRequest{}

			pagination := NewPagination(10, 7, request)
			pagination.offset = 3

			assert.False(t, pagination.CanIncrement())
		},
	)

	t.Run(
		"can increment when offset + limit is less than totalCount",
		func(t *testing.T) {
			request := publicCloud.ApiGetInstanceListRequest{}

			pagination := NewPagination(2, 10, request)
			pagination.offset = 7

			assert.True(t, pagination.CanIncrement())
		},
	)
}

func Test_pagination_nextPage(t *testing.T) {
	t.Run(
		"calling NextPage returns an error if we can't increment",
		func(t *testing.T) {
			request := publicCloud.ApiGetInstanceListRequest{}

			pagination := NewPagination(2, 10, request)
			pagination.offset = 11

			newRequest, err := pagination.NextPage()

			assert.Equal(t, request, newRequest)
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
			request := publicCloud.ApiGetInstanceListRequest{}

			pagination := NewPagination(2, 10, request)
			pagination.offset = 5

			newRequest, err := pagination.NextPage()

			assert.NoError(t, err)
			assert.Equal(t, 7, pagination.offset)

			assert.Equal(
				t,
				int64(7),
				reflect.ValueOf(newRequest).FieldByName("offset").Elem().Int(),
				"request offset is set properly",
			)
		},
	)
}

func Example() {
	request := exampleRequest{CurrentOffset: 0}

	pagination := NewPagination(10, 20, request)

	fmt.Println(request.CurrentOffset)
	fmt.Println(pagination.CanIncrement())
	newRequest, _ := pagination.NextPage()
	fmt.Println(newRequest.CurrentOffset)
	fmt.Println(pagination.CanIncrement())
	newRequest, _ = pagination.NextPage()
	fmt.Println(newRequest.CurrentOffset)
	fmt.Println(pagination.CanIncrement())

	// Output:
	// 0
	// true
	// 10
	// true
	// 20
	// false
}
