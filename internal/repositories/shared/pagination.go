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
	Offset     int
	Limit      int
	TotalCount int
}

type Response interface {
	GetMetadata() publicCloud.Metadata
}

func (p *Pagination) CanIncrement() bool {
	return p.Offset+p.Limit < p.TotalCount
}

func (p *Pagination) NextPage() error {
	if !p.CanIncrement() {
		return ErrCannotIncrementPagination{
			msg: fmt.Sprintf(
				"cannot increment as next offset %d is larger than the total %d",
				p.Offset+p.TotalCount,
				p.TotalCount,
			),
		}
	}

	p.Offset += p.Limit
	return nil
}

func NewPagination(metadata publicCloud.Metadata) Pagination {
	return Pagination{
		Offset:     int(metadata.GetOffset()),
		Limit:      int(metadata.GetLimit()),
		TotalCount: int(metadata.GetTotalCount()),
	}
}
