package service

import (
	"1337bo4rd/internal/core/port"
)

type CommentService struct {
	repo port.CommentRepository
}

func NewCommentService(repo port.CommentRepository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

// func (s *CommentService) GetLastComment(id *uint64) (*domain.Comment, error) {
// 	return s.repo.GetLastComment(id)
// }
