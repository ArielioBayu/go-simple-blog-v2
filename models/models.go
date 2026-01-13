package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	FullName  string    `gorm:"type:varchar(200)" json:"full_name"`
	Bio       string    `gorm:"type:text" json:"bio"`
	RoleID    string    `gorm:"type:char(36);not null" json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
	Posts     []Post    `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
	Comments  []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// Role represents user roles (admin, author, reader)
type Role struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Users       []User    `gorm:"foreignKey:RoleID" json:"users,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

// Post represents a blog post (renamed from Film)
type Post struct {
	ID          string     `gorm:"type:char(36);primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Excerpt     string     `gorm:"type:text" json:"excerpt"`
	AuthorID    string     `gorm:"type:char(36);not null" json:"author_id"`
	Author      User       `gorm:"foreignKey:AuthorID" json:"author"`
	CategoryID  string     `gorm:"type:char(36)" json:"category_id"`
	Category    Category   `gorm:"foreignKey:CategoryID" json:"category"`
	Tags        []Tag      `gorm:"many2many:post_tags;" json:"tags,omitempty"`
	Comments    []Comment  `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Status      string     `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft, published, archived
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// Category represents blog post categories
type Category struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Posts       []Post    `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// Tag represents post tags
type Tag struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
	Posts     []Post    `gorm:"many2many:post_tags;" json:"posts,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// Comment represents comments on posts
type Comment struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	PostID    string    `gorm:"type:char(36);not null" json:"post_id"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post"`
	UserID    string    `gorm:"type:char(36);not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ParentID  *string   `gorm:"type:char(36)" json:"parent_id"` // For nested comments
	Replies   []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
