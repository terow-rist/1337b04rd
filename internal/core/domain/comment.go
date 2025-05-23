package domain

import "time"

type Comment struct {
	ID              uint64
	UserName        string
	UserAvatar      string
	PostID          uint64
	ParentCommentID uint64
	Content         string
	CreatedAt       time.Time
}

type CommentNode struct {
	*Comment
	Replies []*CommentNode
}
