# Changelog

## Version 2.0.0 (2026-01-13)

Complete rewrite and enhancement of the simple blog application from v1.

### 🚀 New Features

#### Authentication & Authorization
- JWT-based authentication system
- User registration and login
- Role-based access control (Admin, Author, Reader)
- Password hashing with bcrypt
- Protected routes with middleware

#### Enhanced Blog Features
- **Posts** (renamed from Film)
  - Title, slug, content, excerpt
  - Post status: draft, published, archived
  - Slug-based URLs
  - Published date tracking
  - Author association
- **Categories**
  - Category management
  - Posts-to-category relationship
  - Admin-only category operations
- **Tags**
  - Many-to-many relationship with posts
  - Automatic slug generation
  - Dynamic tag creation
- **Comments**
  - Nested comments support (replies)
  - User-to-comment association
  - Edit and delete own comments

#### API Enhancements
- RESTful API structure under `/api/v1` prefix
- Pagination support for list endpoints
- Search functionality for posts
- Filtering by category and status
- Proper HTTP status codes
- Consistent error responses

#### Infrastructure
- Environment-based configuration (.env support)
- Structured project layout
- Middleware architecture (Auth, CORS)
- Input validation with binding
- Auto-migration on startup
- Default role seeding
- Docker support with docker-compose
- Makefile for common tasks

#### Documentation
- Comprehensive README with setup instructions
- Detailed API documentation (API.md)
- Docker deployment guide
- Code examples for all endpoints

### 🔐 Security Improvements
- Password hashing (bcrypt)
- JWT token validation
- Input sanitization
- SQL injection protection (GORM)
- Role-based permissions
- CORS configuration

### 🏗️ Architecture Changes
- Modular structure (config, controllers, middleware, models, routes, utils)
- Separation of concerns
- Reusable utilities (JWT, password, slug)
- Better error handling
- Cleaner codebase

### 📊 Database Changes
- `Film` table renamed to `posts` with enhanced fields
- New `categories` table
- New `tags` table
- New `comments` table
- New `post_tags` junction table (many-to-many)
- Enhanced `users` table with email, full_name, bio
- Enhanced `roles` table with description

### 🔄 Breaking Changes from V1

#### API Endpoints
- **Old:** `GET /films` → **New:** `GET /api/v1/posts`
- **Old:** `GET /films/:id` → **New:** `GET /api/v1/posts/:id`
- **Old:** `POST /films` → **New:** `POST /api/v1/posts` (requires auth)
- **Old:** `PUT /films/:id` → **New:** `PUT /api/v1/posts/:id` (requires auth)
- **Old:** `DELETE /films` → **New:** `DELETE /api/v1/posts/:id` (requires auth)

#### Model Changes
- `Film` model replaced with `Post` model
- New fields: `slug`, `excerpt`, `category_id`, `status`, `published_at`
- Removed fields: None (all Film fields mapped to Post)
- Field mappings:
  - `author` (string) → `author_id` (UUID, foreign key)
  - `post` (string) → `content` (text)
  - `category` (string) → `category_id` (UUID, foreign key)

#### Configuration
- Database connection now uses environment variables
- JWT secret configuration required
- CORS origins configurable

### 📈 Performance Improvements
- Regex patterns compiled once at package level
- Efficient slug uniqueness checking
- Proper database indexing (unique constraints on slugs)
- Connection pooling via GORM

### 🧪 Testing & Quality
- Code passes go vet checks
- Code formatted with gofmt
- No security vulnerabilities found (CodeQL scan)
- No dependency vulnerabilities

### 🛠️ Developer Experience
- Makefile for common tasks
- Docker setup for easy deployment
- Environment configuration examples
- Comprehensive documentation
- Clear project structure

### 📝 Migration Guide from V1 to V2

1. **Database:** Create new database or migrate data
   - Export Film data from v1
   - Transform to Post format
   - Import into v2 with proper relationships

2. **API Clients:** Update endpoints
   - Add `/api/v1` prefix to all endpoints
   - Replace `films` with `posts`
   - Add authentication headers
   - Update request/response models

3. **Configuration:**
   - Create `.env` file from `.env.example`
   - Set database credentials
   - Set JWT secret
   - Configure CORS origins

4. **User Management:**
   - Create initial admin user
   - Assign appropriate roles
   - Generate API tokens

### 🎯 Future Enhancements (Not in V2)
- User profile editing
- File upload for images
- Post drafts auto-save
- Email notifications
- Social media sharing
- Analytics dashboard
- Rate limiting
- API versioning
- Swagger/OpenAPI documentation
- Integration tests
- CI/CD pipeline

### 👨‍💻 Credits
**Developer:** Arielio Bayu Aji  
**NIM:** 2103015087  
**Program:** TI  
**Version:** 2.0.0  
**Date:** January 13, 2026
