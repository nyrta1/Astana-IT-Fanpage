package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Author    string    `json:"author" bson:"author"`
	Title     string    `json:"title" json:"title"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Comments struct {
	NewsID   primitive.ObjectID `json:"news_id" bson:"news_id"`
	Comments []CommentData      `json:"comments" bson:"comments"`
}

type Tags struct {
	NewsID primitive.ObjectID `json:"news_id" bson:"news_id"`
	Tags   []TagData          `json:"tags" bson:"tags"`
}

type CommentData struct {
	CommentDataID primitive.ObjectID `json:"comment_id" bson:"_id"`
	Username      string             `json:"user" bson:"user"`
	Content       string             `json:"content" bson:"content"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
}

type TagData struct {
	TagName   string    `json:"tag" bson:"tag"`
	Color     string    `json:"color" bson:"color"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
