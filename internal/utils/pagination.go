package utils

func NewOffset(limit, offset, totalCount int32) *int32 {
	newOffset := offset + limit
	if newOffset >= totalCount {
		return nil
	}

	return &newOffset
}
