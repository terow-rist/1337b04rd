package repository_test

import (
	"1337bo4rd/internal/core/domain"
	"errors"
	"testing"
	"time"
)

// fakeRepo simulates PostRepository for testing without a database
type fakeRepo struct {
	posts    []domain.Post
	comments map[uint64][]domain.Comment
}

func (f *fakeRepo) CreatePost(post *domain.Post) error {
	post.ID = uint64(len(f.posts) + 1)
	f.posts = append(f.posts, *post)
	return nil
}

func (f *fakeRepo) ListPosts() ([]domain.Post, error) {
	if len(f.posts) == 0 {
		return nil, errors.New("no rows in result set")
	}
	return f.posts, nil
}

func (f *fakeRepo) GetPostWithCommentsById(id *uint64) (*domain.PostComents, error) {
	for _, p := range f.posts {
		if p.ID == *id {
			return &domain.PostComents{
				Post:     p,
				Comments: f.comments[*id],
			}, nil
		}
	}
	return nil, errors.New("post not found")
}

func (f *fakeRepo) UpdatePostArchivedAt(postID uint64, archivedAt *time.Time) error {
	for i := range f.posts {
		if f.posts[i].ID == postID {
			f.posts[i].ArchivedAt = *archivedAt
			return nil
		}
	}
	return errors.New("post not found")
}

func TestCreateAndListPosts(t *testing.T) {
	repo := &fakeRepo{}

	post := &domain.Post{
		UserName:   "user1",
		UserAvatar: "avatar.png",
		Title:      "Hello",
		Content:    "World",
		CreatedAt:  time.Now(),
	}

	err := repo.CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	posts, err := repo.ListPosts()
	if err != nil {
		t.Fatalf("ListPosts failed: %v", err)
	}
	if len(posts) != 1 || posts[0].Title != "Hello" {
		t.Errorf("unexpected post data: %+v", posts)
	}
}

func TestGetPostWithCommentsById(t *testing.T) {
	repo := &fakeRepo{
		posts: []domain.Post{
			{ID: 1, Title: "Test"},
		},
		comments: map[uint64][]domain.Comment{
			1: {
				{ID: 1, Content: "First comment"},
			},
		},
	}

	id := uint64(1)
	pc, err := repo.GetPostWithCommentsById(&id)
	if err != nil {
		t.Fatalf("GetPostWithCommentsById failed: %v", err)
	}
	if pc.Post.ID != 1 || len(pc.Comments) != 1 {
		t.Errorf("unexpected result: %+v", pc)
	}
}

func TestUpdatePostArchivedAt(t *testing.T) {
	repo := &fakeRepo{
		posts: []domain.Post{
			{ID: 1, Title: "Old"},
		},
	}

	now := time.Now()
	err := repo.UpdatePostArchivedAt(1, &now)
	if err != nil {
		t.Fatalf("UpdatePostArchivedAt failed: %v", err)
	}
	if !repo.posts[0].ArchivedAt.Equal(now) {
		t.Errorf("ArchivedAt not updated correctly: got %v, want %v", repo.posts[0].ArchivedAt, now)
	}
}
