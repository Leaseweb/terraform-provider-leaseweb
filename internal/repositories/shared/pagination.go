package shared

import (
	"fmt"
)

type ErrCannotIncrementPagination struct {
	msg string
}

func (e ErrCannotIncrementPagination) Error() string {
	return e.msg
}

// Pagination handles pagination for the passed request. The request offset is updated every time NextPage is called.
type Pagination[T Request[T]] struct {
	offset     int
	limit      int
	totalCount int
	Request    T
}

// Request is a contract that Pagination.Request must adhere to so it is supported.
type Request[T any] interface {
	Offset(offset int32) T
}

// CanIncrement returns true if there are any results on the next page.
func (p *Pagination[any]) CanIncrement() bool {
	return p.offset+p.limit <= p.totalCount
}

// NextPage returns an updated Request with a new offset.
func (p *Pagination[Request]) NextPage() (*Request, error) {
	if !p.CanIncrement() {
		return nil, ErrCannotIncrementPagination{
			msg: fmt.Sprintf(
				"cannot increment as next offset %d is larger than the total %d",
				p.offset+p.totalCount,
				p.totalCount,
			),
		}
	}

	p.offset += p.limit
	p.Request = p.Request.Offset(int32(p.offset))

	return &p.Request, nil
}

func NewPagination[T Request[T]](
	limit int32,
	totalCount int32,
	request T,
) Pagination[T] {
	return Pagination[T]{
		offset:     0,
		limit:      int(limit),
		totalCount: int(totalCount),
		Request:    request,
	}
}
