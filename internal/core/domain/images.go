package domain

import (
	"fmt"
)

type ErrImageNotFound struct {
	msg string
}

func (e ErrImageNotFound) Error() string {
	return e.msg
}

type Images []Image

func (i Images) FilterById(id string) (*Image, error) {
	for _, image := range i {
		if image.Id == id {
			return &image, nil
		}
	}

	return nil, ErrImageNotFound{msg: fmt.Sprintf(
		"image with id %q not found",
		id,
	)}
}
