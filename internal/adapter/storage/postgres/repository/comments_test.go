package repository

import (
	"1337bo4rd/internal/core/domain"
	"errors"
	"testing"
	"time"
)

type fakeDB struct {
	lastComment *domain.Comment
}

func (f *fakeDB) GetLastComment(id *uint64) (*domain.Comment, error) {
	if f.lastComment == nil {
		return &domain.Comment{}, errors.New("not found")
	}
	return f.lastComment, nil
}

func (f *fakeDB) CreateComment(comment *domain.Comment) error {
	f.lastComment = comment
	return nil
}

// wrap the fakeDB into the same interface methods as CommentRepository
type fakeRepo struct {
	fakeDB
}

func TestCreateComment(t *testing.T) {
	repo := &fakeRepo{}

	comment := &domain.Comment{
		UserName:   "test_user",
		UserAvatar: "avatar_url",
		PostID:     1,
		Content:    "This is a comment",
		CreatedAt:  time.Now(),
	}

	err := repo.CreateComment(comment)
	if err != nil {
		t.Fatalf("CreateComment returned error: %v", err)
	}

	if repo.lastComment == nil || repo.lastComment.Content != "This is a comment" {
		t.Error("CreateComment did not store the comment correctly")
	}
}

func TestGetLastComment(t *testing.T) {
	repo := &fakeRepo{
		fakeDB{
			lastComment: &domain.Comment{
				ID:         1,
				UserName:   "test_user",
				UserAvatar: "avatar_url",
				PostID:     1,
				Content:    "Last comment",
				CreatedAt:  time.Now(),
			},
		},
	}

	id := uint64(1)
	comment, err := repo.GetLastComment(&id)
	if err != nil {
		t.Fatalf("GetLastComment returned error: %v", err)
	}

	if comment == nil || comment.Content != "Last comment" {
		t.Errorf("Expected comment content 'Last comment', got %v", comment.Content)
	}
}
