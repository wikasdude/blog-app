package model

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type CategoryRequest struct {
	Name string `json:"name" example:"Electronics"`
}
