package response

import "errors"

type Error struct {
	Err string `json:"error"`
}

var (
	ErrInternal    = errors.New("internal error")
	ErrInvalidJSON = errors.New("invalid json passed")
)
