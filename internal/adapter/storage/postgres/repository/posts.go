package repository

import (
	"1337bo4rd/internal/core/domain"
	"database/sql"
	"time"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) CreatePost(post *domain.Post) error {
	var (
		query string
		args  []interface{}
	)
	if post.Image != "" {
		query = `
	INSERT INTO posts (user_name, user_avatar, title, content, image)
	VALUES ($1, $2, $3, $4, $5)
	`
		args = []interface{}{post.UserName, post.UserAvatar, post.Title, post.Content, post.Image}
	} else {
		query = `
	INSERT INTO posts (user_name, user_avatar, title, content)
	VALUES ($1, $2, $3, $4)
	`
		args = []interface{}{post.UserName, post.UserAvatar, post.Title, post.Content}
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) ListPosts() ([]domain.Post, error) {
	var posts []domain.Post
	query := `
	SELECT * FROM posts;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post domain.Post
		var image sql.NullString
		var archivedAt sql.NullTime
		if err := rows.Scan(&post.ID, &post.UserName, &post.UserAvatar, &post.Title, &post.Content, &image, &post.CreatedAt, &archivedAt); err != nil {
			return nil, err
		}

		if image.Valid {
			post.Image = image.String
		}
		if archivedAt.Valid {
			post.ArchivedAt = archivedAt.Time
		}

		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return nil, sql.ErrNoRows
	}

	return posts, nil
}

func (r *PostRepository) GetPostWithCommentsById(id *uint64) (*domain.PostComents, error) {
	query := `
	SELECT 
		p.post_id, p.user_name, p.user_avatar, p.title, p.content, p.image, p.created_at, p.archived_at,
		c.comment_id, c.user_name, c.user_avatar, c.post_id, c.parent_comment_id, c.content, c.created_at
	FROM posts p
	LEFT JOIN comments c ON c.post_id = p.post_id
	WHERE p.post_id = $1
	ORDER BY c.created_at ASC
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post domain.Post
	var comments []domain.Comment
	seen := false

	for rows.Next() {
		var (
			postID       uint64
			postUserName string
			postAvatar   string
			title        string
			content      string
			image        sql.NullString
			createdAt    time.Time
			archivedAt   sql.NullTime

			commentID        sql.NullInt64
			commentUserName  sql.NullString
			commentAvatar    sql.NullString
			commentPostID    sql.NullInt64
			parentCommentID  sql.NullInt64
			commentContent   sql.NullString
			commentCreatedAt sql.NullTime
		)

		if err := rows.Scan(
			&postID, &postUserName, &postAvatar, &title, &content, &image, &createdAt, &archivedAt,
			&commentID, &commentUserName, &commentAvatar, &commentPostID, &parentCommentID, &commentContent, &commentCreatedAt,
		); err != nil {
			return nil, err
		}

		if !seen {
			post = domain.Post{
				ID:         postID,
				UserName:   postUserName,
				UserAvatar: postAvatar,
				Title:      title,
				Content:    content,
				Image:      image.String,
				CreatedAt:  createdAt,
			}

			if image.Valid {
				post.Image = image.String
			}
			if archivedAt.Valid {
				post.ArchivedAt = archivedAt.Time
			}

			seen = true
		}

		if commentID.Valid {
			comment := domain.Comment{
				ID:              uint64(commentID.Int64),
				UserName:        commentUserName.String,
				UserAvatar:      commentAvatar.String,
				PostID:          uint64(commentPostID.Int64),
				ParentCommentID: uint64(parentCommentID.Int64),
				Content:         commentContent.String,
				CreatedAt:       commentCreatedAt.Time,
			}
			comments = append(comments, comment)
		}
	}

	if !seen {
		return nil, sql.ErrNoRows
	}

	return &domain.PostComents{
		Post:     post,
		Comments: comments,
	}, nil
}

func (r *PostRepository) UpdatePostArchivedAt(postID uint64, archivedAt *time.Time) error {
	query := `
		UPDATE posts
		SET archived_at = $1
		WHERE post_id = $2
	`

	_, err := r.db.Exec(query, archivedAt, postID)
	return err
}
