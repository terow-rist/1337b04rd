package service

import (
	"1337bo4rd/internal/core/domain"
	"testing"
	"time"
)

// --- Mock Implementations ---

type mockPostRepo struct {
	createdPost *domain.Post
	createdBy   string
}

func (m *mockPostRepo) CreatePost(post *domain.Post, userId string) error {
	m.createdPost = post
	m.createdBy = userId
	return nil
}

func (m *mockPostRepo) ListPosts() ([]domain.Post, error) {
	return nil, nil
}

func (m *mockPostRepo) UpdatePostArchivedAt(postID uint64, archivedAt *time.Time) error {
	return nil
}

func (m *mockPostRepo) GetPostWithCommentsById(id *uint64) (*domain.PostComents, error) {
	return nil, nil
}

type mockCommentRepo struct{}

func (m *mockCommentRepo) GetLastComment(id *uint64) (*domain.Comment, error) {
	return nil, nil
}

func (m *mockCommentRepo) CreateComment(comment *domain.Comment, userId string) error {
	return nil
}

// --- Test ---

func TestCreatePost(t *testing.T) {
	mockRepo := &mockPostRepo{}
	mockComment := &mockCommentRepo{}
	service := NewPostService(mockRepo, mockComment)

	post := &domain.Post{
		Title:     "Test Title",
		Content:   "Test Content",
		CreatedAt: time.Now(),
	}

	userId := "user-123"
	err := service.CreatePost(post, userId)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	if mockRepo.createdPost != post {
		t.Errorf("Expected post to be created")
	}
	if mockRepo.createdBy != userId {
		t.Errorf("Expected userId %q; got %q", userId, mockRepo.createdBy)
	}
}
