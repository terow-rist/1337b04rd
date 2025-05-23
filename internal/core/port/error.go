package port

import "errors"

var (
	ErrNoPosts       = errors.New("no posts")
	ErrInvalidPostId = errors.New("invalid post ID")
)
