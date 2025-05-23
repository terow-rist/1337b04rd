package service

import (
	"1337bo4rd/internal/core/domain"
	"1337bo4rd/internal/core/port"
	"database/sql"
	"errors"
	"strconv"
	"time"
)

type PostService struct {
	repo        port.PostRepository
	commentRepo port.CommentRepository
}

func NewPostService(repo port.PostRepository, commentRepo port.CommentRepository) *PostService {
	return &PostService{
		repo:        repo,
		commentRepo: commentRepo,
	}
}

func (s *PostService) CreatePost(post *domain.Post) error {
	return s.repo.CreatePost(post)
}

func (s *PostService) ListPosts() ([]domain.Post, error) {
	posts, err := s.repo.ListPosts()
	if err != nil {
		return nil, err
	}

	return posts, nil

}

func (s *PostService) ListActive() ([]domain.Post, error) {
	posts, err := s.repo.ListPosts()
	if err != nil {
		return nil, err
	}

	var validPosts []domain.Post
	now := time.Now()

	for _, post := range posts {
		comment, err := s.commentRepo.GetLastComment(&post.ID)
		if err != nil {
			// No comment found, fallback to post time
			if errors.Is(err, sql.ErrNoRows) {
				if now.Sub(post.CreatedAt) < 10*time.Minute {
					validPosts = append(validPosts, post)
				} else {
					err = s.repo.UpdatePostArchivedAt(post.ID, &now)
					if err != nil {
						return nil, err
					}
				}
				continue
			}
			// Some other error
			return nil, err
		}

		// If comment exists, check if it's recent enough
		if now.Sub(comment.CreatedAt) < 15*time.Minute {
			validPosts = append(validPosts, post)
		} else {
			err = s.repo.UpdatePostArchivedAt(post.ID, &now)
			if err != nil {
				return nil, err
			}
		}
	}

	return validPosts, nil
}

func (s *PostService) GetPostWithCommentsById(idStr string) (*domain.PostComents, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, port.ErrInvalidPostId
	}
	return s.repo.GetPostWithCommentsById(&id)
}

func (s *PostService) CreateComment(comment *domain.Comment, idStr string) error {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return port.ErrInvalidPostId
	}

	comment.PostID = id

	if err := s.commentRepo.CreateComment(comment); err != nil {
		return err
	}
	return nil
}
