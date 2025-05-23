package domain

import "time"

type Post struct {
	ID         uint64
	UserName   string
	UserAvatar string
	Title      string
	Content    string
	Image      string
	CreatedAt  time.Time
	ArchivedAt time.Time
}

type PostComents struct {
	Post     Post
	Comments []Comment
}
