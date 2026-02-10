package memberships

import "time"

type SignUpRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

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
