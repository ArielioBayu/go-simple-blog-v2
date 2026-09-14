package user

import (
	"time"
)

type UserModel struct {
	ID        int64     `column:"id"`
	Email     string    `column:"email"`
	Password  string    `column:"password"`
	CreatedAt time.Time `column:"created_at"`
	UpdatedAt time.Time `column:"updated_at"`
	CreatedBy string    `column:"created_by"`
	UpdatedBy string    `column:"updated_by"`
	Username  string    `column:"username"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
