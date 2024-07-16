package value_object

import (
	"fmt"

	"github.com/google/uuid"
)

type errCannotConvertValueToUUID struct {
	msg string
}

func (e errCannotConvertValueToUUID) Error() string {
	return e.msg
}

type Uuid struct {
	Uuid uuid.UUID
}

func (u Uuid) String() string {
	return u.Uuid.String()
}

func NewUuid(value string) (*Uuid, error) {
	parsedUuid, err := uuid.Parse(value)
	if err != nil {
		return nil, errCannotConvertValueToUUID{msg: fmt.Sprintf("could not convert %q to UUID", value)}
	}

	return &Uuid{Uuid: parsedUuid}, nil
}

func NewGeneratedUuid() Uuid {
	return Uuid{Uuid: uuid.New()}
}
