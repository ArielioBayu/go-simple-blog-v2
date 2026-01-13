# API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Most endpoints require a JWT token. Include it in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

---

## Authentication Endpoints

### Register
Create a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

**Response (201):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "role": {
      "id": "uuid",
      "name": "reader",
      "description": "Can read and comment on posts"
    }
  }
}
```

### Login
Authenticate a user.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "role": {
      "id": "uuid",
      "name": "reader"
    }
  }
}
```

### Get Current User
Get authenticated user information.

**Endpoint:** `GET /auth/me` 🔒

**Response (200):**
```json
{
  "user": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "bio": "",
    "role": {
      "id": "uuid",
      "name": "reader"
    }
  }
}
```

---

## Post Endpoints

### List Posts
Get all published posts with pagination and search.

**Endpoint:** `GET /posts`

**Query Parameters:**
- `page` (int): Page number (default: 1)
- `page_size` (int): Items per page (default: 10, max: 100)
- `search` (string): Search in title and content
- `category_id` (string): Filter by category
- `status` (string): Filter by status (default: "published")

**Example:** `GET /posts?page=1&page_size=10&search=golang`

**Response (200):**
```json
{
  "posts": [
    {
      "id": "uuid",
      "title": "My First Post",
      "slug": "my-first-post",
      "content": "Content here...",
      "excerpt": "Brief summary",
      "author_id": "uuid",
      "author": {
        "id": "uuid",
        "username": "john_doe",
        "full_name": "John Doe"
      },
      "category_id": "uuid",
      "category": {
        "id": "uuid",
        "name": "Technology",
        "slug": "technology"
      },
      "tags": [
        {
          "id": "uuid",
          "name": "Golang",
          "slug": "golang"
        }
      ],
      "status": "published",
      "published_at": "2024-01-13T12:00:00Z",
      "created_at": "2024-01-13T12:00:00Z",
      "updated_at": "2024-01-13T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

### Get Single Post
Get a post by ID or slug.

**Endpoint:** `GET /posts/:id`

**Example:** `GET /posts/my-first-post` or `GET /posts/uuid`

**Response (200):**
```json
{
  "post": {
    "id": "uuid",
    "title": "My First Post",
    "slug": "my-first-post",
    "content": "Full content...",
    "excerpt": "Brief summary",
    "author": { },
    "category": { },
    "tags": [ ],
    "comments": [ ],
    "status": "published",
    "published_at": "2024-01-13T12:00:00Z",
    "created_at": "2024-01-13T12:00:00Z",
    "updated_at": "2024-01-13T12:00:00Z"
  }
}
```

### Get My Posts
Get posts created by authenticated user.

**Endpoint:** `GET /posts/my-posts` 🔒

**Query Parameters:**
- `page` (int): Page number
- `page_size` (int): Items per page

**Response (200):** Same format as List Posts

### Create Post
Create a new post (requires author or admin role).

**Endpoint:** `POST /posts` 🔒

**Request Body:**
```json
{
  "title": "My New Post",
  "content": "Full content of the post...",
  "excerpt": "Brief summary",
  "category_id": "uuid",
  "tags": ["golang", "api"],
  "status": "published"
}
```

**Response (201):**
```json
{
  "message": "Post created successfully",
  "post": { }
}
```

### Update Post
Update a post (requires author/admin).

**Endpoint:** `PUT /posts/:id` 🔒

**Request Body:**
```json
{
  "title": "Updated Title",
  "content": "Updated content...",
  "status": "published"
}
```

**Response (200):**
```json
{
  "message": "Post updated successfully",
  "post": { }
}
```

### Delete Post
Delete a post (requires author/admin).

**Endpoint:** `DELETE /posts/:id` 🔒

**Response (200):**
```json
{
  "message": "Post deleted successfully"
}
```

---

## Category Endpoints

### List Categories
Get all categories.

**Endpoint:** `GET /categories`

**Response (200):**
```json
{
  "categories": [
    {
      "id": "uuid",
      "name": "Technology",
      "slug": "technology",
      "description": "Tech-related posts",
      "created_at": "2024-01-13T12:00:00Z",
      "updated_at": "2024-01-13T12:00:00Z"
    }
  ]
}
```

### Get Single Category
Get a category with its posts.

**Endpoint:** `GET /categories/:id`

**Response (200):**
```json
{
  "category": {
    "id": "uuid",
    "name": "Technology",
    "slug": "technology",
    "description": "Tech-related posts",
    "posts": [ ]
  }
}
```

### Create Category
Create a new category (admin only).

**Endpoint:** `POST /categories` 🔒 👑

**Request Body:**
```json
{
  "name": "Technology",
  "description": "Tech-related posts"
}
```

**Response (201):**
```json
{
  "message": "Category created successfully",
  "category": { }
}
```

### Update Category
Update a category (admin only).

**Endpoint:** `PUT /categories/:id` 🔒 👑

**Request Body:**
```json
{
  "name": "Updated Name",
  "description": "Updated description"
}
```

**Response (200):**
```json
{
  "message": "Category updated successfully",
  "category": { }
}
```

### Delete Category
Delete a category (admin only).

**Endpoint:** `DELETE /categories/:id` 🔒 👑

**Response (200):**
```json
{
  "message": "Category deleted successfully"
}
```

---

## Comment Endpoints

### List Comments
Get comments for a post.

**Endpoint:** `GET /comments?post_id=uuid`

**Response (200):**
```json
{
  "comments": [
    {
      "id": "uuid",
      "content": "Great post!",
      "post_id": "uuid",
      "user_id": "uuid",
      "user": {
        "id": "uuid",
        "username": "john_doe"
      },
      "parent_id": null,
      "replies": [ ],
      "created_at": "2024-01-13T12:00:00Z",
      "updated_at": "2024-01-13T12:00:00Z"
    }
  ]
}
```

### Create Comment
Create a comment or reply.

**Endpoint:** `POST /comments` 🔒

**Request Body:**
```json
{
  "content": "Great post!",
  "post_id": "uuid",
  "parent_id": "uuid"
}
```

**Response (201):**
```json
{
  "message": "Comment created successfully",
  "comment": { }
}
```

### Update Comment
Update a comment.

**Endpoint:** `PUT /comments/:id` 🔒

**Request Body:**
```json
{
  "content": "Updated comment"
}
```

**Response (200):**
```json
{
  "message": "Comment updated successfully",
  "comment": { }
}
```

### Delete Comment
Delete a comment.

**Endpoint:** `DELETE /comments/:id` 🔒

**Response (200):**
```json
{
  "message": "Comment deleted successfully"
}
```

---

## Error Responses

All endpoints may return error responses:

**400 Bad Request:**
```json
{
  "error": "Validation error message"
}
```

**401 Unauthorized:**
```json
{
  "error": "Authorization header is required"
}
```

**403 Forbidden:**
```json
{
  "error": "Insufficient permissions"
}
```

**404 Not Found:**
```json
{
  "error": "Resource not found"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Internal server error"
}
```

---

## Icons Legend
- 🔒 Requires authentication
- 👑 Requires admin role
