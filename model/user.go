package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // hide password in JSON
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
