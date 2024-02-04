package models

import (
	"time"
)

// PostgresSQL Models

type User struct {
	ID        uint      `json:"ID,omitempty"`
	CreatedAt time.Time `json:"CreatedAt"`
	Username  string    `json:"username" json:"Username,omitempty"`
	Password  string    `json:"password" json:"Password,omitempty"`
	UserType  uint      `json:"userType"`
}

// MongoDB Models

type News struct {
	Tags      []Tag     `json:"tags" bson:"tags"`
	Author    string    `json:"author" bson:"author"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Comments  []Comment `json:"comments" bson:"comments"`
}

type Comment struct {
	Username  string    `json:"user" bson:"user"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Tag struct {
	Tag       string    `json:"tag" bson:"tag"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
