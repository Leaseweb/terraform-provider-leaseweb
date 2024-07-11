package value_object

import (
	"errors"

	"github.com/google/uuid"
)

var ErrCouldNotConvertValueIntoUUID = errors.New("could not convert Uuid into UUID")

type Uuid struct {
	Uuid uuid.UUID
}

func (u Uuid) String() string {
	return u.Uuid.String()
}

func NewUuid(value string) (*Uuid, error) {
	parsedUuid, err := uuid.Parse(value)
	if err != nil {
		return nil, ErrCouldNotConvertValueIntoUUID
	}

	return &Uuid{Uuid: parsedUuid}, nil
}

func NewGeneratedUuid() Uuid {
	return Uuid{Uuid: uuid.New()}
}
