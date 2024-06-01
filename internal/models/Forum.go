package models

import (
	"time"
)

type Forum struct {
	ID            int        `json:"id"`
	Forum         string     `json:"forum_text"`
	UserID        int        `json:"user_id"`
	PublishedAt   time.Time  `json:"published_at"`
	Comments      []*Comment `json:"comments,omitempty"`
	Likes         []*Like    `json:"likes,omitempty"`
	CommentsArray []int      `json:"comments_array,omitempty"`
	LikesArray    []int      `json:"likes_array,omitempty"`
}

type Comment struct {
	ID          int       `json:"id"`
	ForumID     int       `json:"forum_id"`
	UserID      int       `json:"user_id"`
	Comment     string    `json:"comment_text"`
	PublishedAt time.Time `json:"published_at"`
}

type Like struct {
	ID      int `json:"id"`
	ForumID int `json:"forum_id"`
	UserID  int `json:"user_id"`
}
