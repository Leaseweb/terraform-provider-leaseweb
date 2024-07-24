package shared

import (
	"fmt"

	"github.com/leaseweb/leaseweb-go-sdk/publicCloud"
)

type ErrCannotIncrementPagination struct {
	msg string
}

func (e ErrCannotIncrementPagination) Error() string {
	return e.msg
}

type Pagination struct {
	offset     int
	limit      int
	totalCount int
	Request    publicCloud.ApiGetInstanceListRequest
}

type Response interface {
	GetMetadata() publicCloud.Metadata
}

func (p *Pagination) CanIncrement() bool {
	return p.offset+p.limit < p.totalCount
}

func (p *Pagination) NextPage() (*publicCloud.ApiGetInstanceListRequest, error) {
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

func NewPagination(
	limit int32,
	totalCount int32,
	request publicCloud.ApiGetInstanceListRequest,
) Pagination {
	return Pagination{
		offset:     0,
		limit:      int(limit),
		totalCount: int(totalCount),
		Request:    request,
	}
}
