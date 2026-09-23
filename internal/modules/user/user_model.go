package user

import (
	"time"
)

type UserModel struct {
	ID         int64     `column:"id"`
	Email      string    `column:"email"`
	Bio        string    `column:"bio"`
	AvatarURL  string    `column:"avatar_url"`
	BannerURL  string    `column:"banner_url"`
	IsVerified bool      `column:"is_verified"`
	Password   string    `column:"password"`
	CreatedAt  time.Time `column:"created_at"`
	UpdatedAt  time.Time `column:"updated_at"`
	CreatedBy  string    `column:"created_by"`
	UpdatedBy  string    `column:"updated_by"`
	Username   string    `column:"username"`
}

type UserResponse struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Bio        string    `json:"bio"`
	AvatarURL  string    `json:"avatar_url"`
	BannerURL  string    `json:"banner_url"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

type ProfileStats struct {
	StoriesCount   int `json:"stories_count"`
	SavedCount     int `json:"saved_count"`
	LikesCount     int `json:"likes_count"`
	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
}

type ProfileResponse struct {
	ID         int64        `json:"id"`
	Username   string       `json:"username"`
	Email      string       `json:"email"`
	Bio        string       `json:"bio"`
	AvatarURL  string       `json:"avatar_url"`
	BannerURL  string       `json:"banner_url"`
	IsVerified bool         `json:"is_verified"`
	CreatedAt  time.Time    `json:"created_at"`
	Stats      ProfileStats `json:"stats"`
}

type UpdateProfileRequest struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	BannerURL string `json:"banner_url"`
}
