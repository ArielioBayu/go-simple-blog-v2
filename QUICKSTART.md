# Quick Start Guide

Get the simple blog v2 up and running in minutes!

## Method 1: Using Docker (Recommended)

The fastest way to get started:

```bash
# Clone the repository
git clone https://github.com/ArielioBayu/go-simple-blog-v2.git
cd go-simple-blog-v2

# Start everything with Docker
docker-compose up -d

# Check if it's running
curl http://localhost:8080/api/v1/health
```

That's it! The API is now running on `http://localhost:8080`

## Method 2: Manual Setup

### Prerequisites
- Go 1.22+
- MySQL 5.7+

### Steps

1. **Clone the repository**
```bash
git clone https://github.com/ArielioBayu/go-simple-blog-v2.git
cd go-simple-blog-v2
```

2. **Setup database**
```bash
mysql -u root -p
CREATE DATABASE simple_blog_v2;
exit
```

3. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your settings
```

4. **Run the application**
```bash
make install  # Install dependencies
make build    # Build the binary
make run      # Run the application
```

## First Steps with the API

### 1. Register a user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123",
    "full_name": "Admin User"
  }'
```

### 2. Login to get token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

Copy the `token` from the response.

### 3. Create a post
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "My First Post",
    "content": "This is my first blog post!",
    "excerpt": "An introduction to my blog",
    "status": "published"
  }'
```

### 4. List all posts
```bash
curl http://localhost:8080/api/v1/posts
```

## Common Tasks

### Build the project
```bash
make build
```

### Run tests
```bash
make test
```

### Format code
```bash
make fmt
```

### Clean build artifacts
```bash
make clean
```

### View logs (Docker)
```bash
docker-compose logs -f app
```

### Stop Docker services
```bash
docker-compose down
```

## Next Steps

- Check [README.md](README.md) for complete documentation
- See [API.md](API.md) for full API reference
- Read [CHANGELOG.md](CHANGELOG.md) for migration guide from v1

## Troubleshooting

### Database connection failed
- Check MySQL is running
- Verify credentials in `.env`
- Check database exists

### Port already in use
- Change `PORT` in `.env`
- Or stop the service using port 8080

### Permission denied
```bash
chmod +x simple-blog-v2
```

## Need Help?

Check the documentation:
- [README.md](README.md) - Full documentation
- [API.md](API.md) - API reference
- [CHANGELOG.md](CHANGELOG.md) - Version history

Happy blogging! 🎉
