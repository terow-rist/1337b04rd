package port

import "1337bo4rd/internal/core/domain"

type CommentRepository interface {
	GetLastComment(id *uint64) (*domain.Comment, error)
	CreateComment(comment *domain.Comment) error
}

type CommentService interface {
}
