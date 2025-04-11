package model

import "time"

type Post struct {
	ID          int       `json:"post_id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        *string   `json:"body"`
	CategoryID  int       `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type CreatePostRequest struct {
	Title       string `json:"title" example:"How to Build a hi"`
	Description string `json:"description" example:"This post explains how to make microservices."`
	CategoryID  int    `json:"category_id" example:"2"`
	Body        string `json:"body" example:"hello"`
}
type UpdatePostRequest struct {
	Title       string `json:"title" example:"How to Build a"`
	Description string `json:"description" example:"This post explains how to design and build a RESTful API using Golang."`
	CategoryID  int    `json:"category_id" example:"2"`
	Body        string `json:"body" example:"hi"`
}
