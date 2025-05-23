package repository

import (
	"1337bo4rd/internal/core/domain"
	"database/sql"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) GetLastComment(id *uint64) (*domain.Comment, error) {
	query := `
	SELECT * FROM comments
	WHERE post_id = $1
	ORDER BY comment_id DESC
	LIMIT 1
	`
	rows := r.db.QueryRow(query, id)
	var comment domain.Comment
	var postID sql.NullInt64
	var pcID sql.NullInt64
	if err := rows.Scan(&comment.ID, &comment.UserName, &comment.UserAvatar, &postID, &pcID, &comment.Content, &comment.CreatedAt); err != nil {
		return &domain.Comment{}, err
	}

	return &comment, nil
}

func (r *CommentRepository) CreateComment(comment *domain.Comment) error {
	var query string
	var args []interface{}
	if comment.ParentCommentID == 0 {
		query = `
		INSERT INTO comments (user_name, user_avatar, post_id,content)
		VALUES ($1, $2, $3, $4)
		`
		args = []interface{}{comment.UserName, comment.UserAvatar, comment.PostID, comment.Content}
	} else {
		query = `
		INSERT INTO comments (user_name, user_avatar, post_id, parent_comment_id, content)
		VALUES ($1, $2, $3, $4, $5)
		`
		args = []interface{}{comment.UserName, comment.UserAvatar, comment.PostID, comment.ParentCommentID, comment.Content}

	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
