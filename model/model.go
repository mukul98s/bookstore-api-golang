package model

import "time"

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name" validate:"required,min=3"`
	Email     string    `json:"email" validate:"required,email"`
	Phone     string    `json:"phone" validate:"required,min=10,max=10"`
	Password  string    `json:"password" validate:"required,min=8"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Book struct {
	Id         string    `json:"id"`
	UserID     string    `json:"user_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	AuthorName string    `json:"author_name" validate:"required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}
