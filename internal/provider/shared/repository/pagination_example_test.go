package repository

import (
	"fmt"
)

type exampleRequest struct {
	CurrentOffset int32
}

func (e exampleRequest) Offset(offset int32) exampleRequest {
	e.CurrentOffset = offset
	return e
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
