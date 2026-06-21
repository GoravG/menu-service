package repository

import (
	"errors"
	"strings"
)

var ErrDuplicate = errors.New("duplicate entry")

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "Duplicate entry") {
		return ErrDuplicate
	}
	return err
}
