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

type Pagination[T Request[T]] struct {
	offset     int
	limit      int
	totalCount int
	Request    T
}

type Request[T any] interface {
	Offset(offset int32) T
}

func (p *Pagination[any]) CanIncrement() bool {
	return p.offset+p.limit < p.totalCount
}

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
