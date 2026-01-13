# Go Simple Blog V2

A simple but comprehensive blog API built with Go, enhanced from v1 with modern features and best practices.

## 🚀 Features

### From V1
- ✅ RESTful API architecture
- ✅ Gin web framework
- ✅ GORM ORM
- ✅ MySQL database

### New in V2
- 🔐 **Authentication & Authorization**
  - JWT-based authentication
  - Role-based access control (Admin, Author, Reader)
  - Password hashing with bcrypt
  
- 📝 **Enhanced Blog Features**
  - Post management (renamed from Film)
  - Categories with slug support
  - Tags system (many-to-many relationship)
  - Comments with nested replies support
  - Post status (draft, published, archived)
  
- 🔧 **Better Architecture**
  - Environment-based configuration
  - Structured project layout
  - Middleware (Auth, CORS)
  - Input validation
  - Proper error handling
  
- 📄 **API Features**
  - Pagination support
  - Search functionality
  - Slug-based URLs
  - User profile management

## 📁 Project Structure

```
go-simple-blog-v2/
├── config/              # Configuration and database setup
├── controllers/         # Request handlers
│   ├── authcontroller/
│   ├── postcontroller/
│   ├── categorycontroller/
│   └── commentcontroller/
├── middleware/          # Auth, CORS, etc.
├── models/             # Database models
├── routes/             # Route definitions
├── utils/              # Helper functions (JWT, password, slug)
├── .env.example        # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
└── main.go            # Application entry point
```

## 🛠️ Setup

### Prerequisites
- Go 1.22 or higher
- MySQL 5.7 or higher

### Installation

1. Clone the repository:
```bash
git clone https://github.com/ArielioBayu/go-simple-blog-v2.git
cd go-simple-blog-v2
```

2. Copy `.env.example` to `.env` and configure:
```bash
cp .env.example .env
```

3. Edit `.env` with your database credentials:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=simple_blog_v2
JWT_SECRET=your-secret-key-here
```

4. Create the database:
```sql
CREATE DATABASE simple_blog_v2;
```

5. Install dependencies:
```bash
go mod download
```

6. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## 📚 API Endpoints

### Health Check
```
GET /api/v1/health
```

### Authentication
```
POST   /api/v1/auth/register   - Register new user
POST   /api/v1/auth/login      - Login user
GET    /api/v1/auth/me         - Get current user (protected)
```

### Posts
```
GET    /api/v1/posts           - List all published posts (with pagination & search)
GET    /api/v1/posts/:id       - Get single post (by ID or slug)
GET    /api/v1/posts/my-posts  - Get current user's posts (protected)
POST   /api/v1/posts           - Create new post (protected, author/admin)
PUT    /api/v1/posts/:id       - Update post (protected, owner/admin)
DELETE /api/v1/posts/:id       - Delete post (protected, owner/admin)
```

### Categories
```
GET    /api/v1/categories      - List all categories
GET    /api/v1/categories/:id  - Get single category with posts
POST   /api/v1/categories      - Create category (protected, admin only)
PUT    /api/v1/categories/:id  - Update category (protected, admin only)
DELETE /api/v1/categories/:id  - Delete category (protected, admin only)
```

### Comments
```
GET    /api/v1/comments?post_id=xxx  - List comments for a post
POST   /api/v1/comments              - Create comment (protected)
PUT    /api/v1/comments/:id          - Update comment (protected, owner/admin)
DELETE /api/v1/comments/:id          - Delete comment (protected, owner/admin)
```

## 🔒 Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## 📝 Example Usage

### Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "password123"
  }'
```

### Create a post
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post",
    "excerpt": "A brief summary",
    "status": "published",
    "tags": ["golang", "api"]
  }'
```

### List posts with pagination
```bash
curl "http://localhost:8080/api/v1/posts?page=1&page_size=10&search=golang"
```

## 👥 User Roles

- **Admin**: Full access to all resources
- **Author**: Can create and manage own posts, create comments
- **Reader**: Can read posts and create comments

## 🔐 Security Features

- Password hashing with bcrypt
- JWT token authentication
- Role-based access control
- Input validation
- SQL injection protection (via GORM)

## 🚧 Migration from V1

Key changes from V1:
- `Film` model renamed to `Post` with enhanced fields
- Added `Category`, `Tag`, and `Comment` models
- Added authentication and authorization
- API endpoints now under `/api/v1` prefix
- Environment-based configuration
- Improved error handling and validation

## 👨‍💻 Author

**Arielio Bayu Aji**
- NIM: 2103015087
- Prodi: TI

## 📄 License

This project is for educational purposes.