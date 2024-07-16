package models

import (
	"time"
)

type Forum struct {
	ID             int        `json:"id"`
	Forum          string     `json:"forum_text"`
	UserID         int        `json:"user_id"`
	PublishedAt    time.Time  `json:"published_at"`
	Comments       []*Comment `json:"comments,omitempty"`
	Likes          []*Like    `json:"likes,omitempty"`
	UserUsername   string     `json:"user_username,omitempty"`
	UserName       string     `json:"user_name,omitempty"`
	UserPhoto      string     `json:"user_photo,omitempty"`
	CommentsNumber int        `json:"comments_number,omitempty"`
	LikesNumber    int        `json:"likes_number,omitempty"`
}

type Comment struct {
	ID           int       `json:"id"`
	ForumID      int       `json:"forum_id"`
	UserID       int       `json:"user_id"`
	Comment      string    `json:"reply_text"`
	PublishedAt  time.Time `json:"published_at"`
	UserUsername string    `json:"user_username,omitempty"`
	UserPhoto    string    `json:"user_photo,omitempty"`
	UserName     string    `json:"user_name,omitempty"`
}

type Like struct {
	ID        int       `json:"id"`
	ForumID   int       `json:"forum_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
