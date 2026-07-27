# Dokumentasi Backend — go-simple-blog-v2

> **Versi Dokumentasi:** 1.0.0  
> **Tanggal dibuat:** 2026-07-23  
> **Module:** `github.com/ArielioBayu/go-simple-blog-v2`  
> **Go Version:** 1.24.0

---

## Daftar Isi

1. [Gambaran Umum](#1-gambaran-umum)
2. [Tech Stack](#2-tech-stack)
3. [Struktur Project](#3-struktur-project)
4. [Arsitektur Layer](#4-arsitektur-layer)
5. [Konfigurasi](#5-konfigurasi)
6. [Database & Migrasi](#6-database--migrasi)
7. [Alur Inisialisasi (main.go)](#7-alur-inisialisasi-maingo)
8. [Middleware](#8-middleware)
9. [Fitur: Memberships (Autentikasi)](#9-fitur-memberships-autentikasi)
10. [Fitur: Posts (Blog)](#10-fitur-posts-blog)
11. [Package Shared (pkg/)](#11-package-shared-pkg)
12. [Error Constants](#12-error-constants)
13. [Alur Request End-to-End](#13-alur-request-end-to-end)
14. [Ringkasan Endpoint API](#14-ringkasan-endpoint-api)
15. [Catatan Penting & Perhatian](#15-catatan-penting--perhatian)
16. [Changelog Dokumentasi](#16-changelog-dokumentasi)

---

## 1. Gambaran Umum

`go-simple-blog-v2` adalah RESTful API backend untuk platform blog sederhana yang dibangun menggunakan bahasa Go. Project ini menerapkan arsitektur berlapis (layered architecture) dengan pola **Handler → Service → Repository → Database**.

**Fitur utama:**
- Registrasi dan Login pengguna dengan password terenkripsi (bcrypt)
- Autentikasi berbasis JWT (Access Token) + Refresh Token
- CRUD Posts (artikel blog) dengan dukungan hashtag
- Sistem komentar pada post
- Sistem aktivitas (like/unlike) pada post
- Pagination pada daftar post

---

## 2. Tech Stack

| Komponen       | Library / Tool                      | Versi     |
|----------------|-------------------------------------|-----------|
| Web Framework  | `github.com/gin-gonic/gin`          | v1.11.0   |
| Database       | MySQL                               | -         |
| DB Driver      | `github.com/go-sql-driver/mysql`    | v1.9.3    |
| JWT            | `github.com/golang-jwt/jwt/v5`      | v5.3.1    |
| Config Loader  | `github.com/spf13/viper`            | v1.21.0   |
| Enkripsi       | `golang.org/x/crypto/bcrypt`        | v0.40.0   |
| Migration Tool | `golang-migrate`                    | (CLI)     |

---

## 3. Struktur Project

```
go-simple-blog-v2/
├── main.go                              # Entry point aplikasi
├── go.mod / go.sum                      # Dependency management
├── Makefile                             # Perintah migrate
├── DOCUMENTATION.md                     # Dokumentasi project ini
│
├── internal/                            # Kode inti aplikasi (tidak di-export ke luar)
│   ├── configs/
│   │   ├── config.go                    # Logic load config menggunakan Viper
│   │   ├── types.go                     # Struct definisi konfigurasi
│   │   └── config.yaml.example          # Contoh file konfigurasi
│   │
│   ├── constants/
│   │   └── errors.go                    # Definisi semua error sentral
│   │
│   ├── middleware/
│   │   ├── cors.go                      # Middleware CORS
│   │   └── middleware.go                # Auth middleware (Header & Cookie)
│   │
│   ├── model/
│   │   ├── memberships/
│   │   │   └── memberships_model.go     # Model request/response & DB untuk user
│   │   └── posts/
│   │       ├── post_model.go            # Model request/response & DB untuk post
│   │       ├── comment_model.go         # Model request/response & DB untuk komentar
│   │       └── actvities_model.go       # Model request/response & DB untuk aktivitas
│   │
│   ├── handlers/                        # Layer HTTP Handler
│   │   ├── memberships/
│   │   │   ├── memberships_handler.go   # Inisialisasi handler & registrasi route
│   │   │   ├── login.go                 # Handler SignUp & SignIn
│   │   │   ├── get_user.go              # Handler GetUser
│   │   │   └── refresh_token.go         # Handler Refresh Token
│   │   └── posts/
│   │       ├── posts_handler.go         # Inisialisasi handler & registrasi route
│   │       ├── post.go                  # Handler CreatePost, GetAllPost, GetPostById
│   │       ├── comment.go               # Handler CreateComment
│   │       └── activities.go            # Handler InsertUpdateActivities
│   │
│   ├── service/                         # Layer Business Logic
│   │   ├── memberships/
│   │   │   ├── service.go               # Interface & constructor MembershipsService
│   │   │   ├── login.go                 # Logic GetUser, SignUp, SignIn
│   │   │   └── refresh_token.go         # Logic GetIdRefreshToken, ValidateRefreshToken
│   │   └── posts/
│   │       ├── service.go               # Interface & constructor PostsService
│   │       ├── post.go                  # Logic CreatePost, GetAllPost, GetPostById
│   │       ├── comment.go               # Logic CreateComment
│   │       └── activities.go            # Logic InsertUpdateActivities
│   │
│   └── repository/                      # Layer Akses Database (raw SQL)
│       ├── memberships/
│       │   ├── repository.go            # Interface & constructor MembershipRepository
│       │   ├── users.go                 # Query: GetUser, GetUserById, GetUserByEmail, CreateUser
│       │   └── refresh_token.go         # Query: InsertRefreshToken, GetRefreshToken, GetIdRefreshToken
│       └── posts/
│           ├── repository.go            # Interface & constructor PostsRepository
│           ├── post.go                  # Query: CreatePost, GetAllPost, GetPostById
│           ├── comment.go               # Query: GetCommentById, CreateComment
│           └── activities.go            # Query: CountLikedByPostID, GetActivities, CreateActivities, UpdateActivities
│
├── pkg/                                 # Package shared yang bisa dipakai lintas domain
│   ├── internalsql/
│   │   └── sql.go                       # Fungsi koneksi ke MySQL
│   ├── jwt/
│   │   └── jwt.go                       # Fungsi CreateToken & ValidateToken
│   ├── token/
│   │   └── refresh_token.go             # Fungsi GenerateRefreshToken (random hex)
│   └── utils/
│       └── helper.go                    # JsonTime, SetAccessTokenCookie, SetCookie
│
└── scripts/
    └── migrations/                      # File SQL migrasi database (urut berdasarkan nomor)
        ├── 000001_create_users_table
        ├── 000002_alter_create_username
        ├── 000003_create_post_table
        ├── 000004_alter_foreign_key_post_table
        ├── 000005_create_comment_table
        ├── 000006_create_activities_table
        └── 000007_create_refresh_tokens_table
```

---

## 4. Arsitektur Layer

Project menggunakan **Clean Architecture / Layered Architecture** dengan arah dependensi dari luar ke dalam:

```
Request HTTP
     ↓
[ Handler ]   ← menerima & validasi request, format response
     ↓
[ Service ]   ← business logic, orchestrasi, error handling
     ↓
[Repository]  ← raw SQL query ke database
     ↓
[ Database ]  ← MySQL
```

**Prinsip penting:**
- Setiap layer hanya berkomunikasi dengan layer di bawahnya melalui **interface**
- Handler tidak boleh langsung memanggil repository
- Repository tidak mengetahui apa-apa tentang business logic
- Model/struct dipisah per domain (`memberships` / `posts`)

---

## 5. Konfigurasi

### File: `internal/configs/`

Konfigurasi dibaca dari file `config.yaml` menggunakan **Viper**. Semua setting dimasukkan ke struct `Config`.

```yaml
# config.yaml (buat dari config.yaml.example)
service:
  port: ":9888"
  secret_key: "your_secret_key_here"

database:
  dbsourcename: "username:password@tcp(localhost:3306)/database_name?parseTime=true&loc=Asia%2FJakarta"
```

### Struct Konfigurasi

```go
type Config struct {
    Service  Service  // port & secret_key
    Database Database // dsn string
}
```

### Cara Kerja Init Config

1. `configs.Init()` dipanggil di `main.go` dengan opsi path folder, nama file, dan tipe file
2. Viper membaca file yaml, mendukung override lewat environment variable (`AutomaticEnv`)
3. Hasil unmarshal disimpan di variabel global `config *Config`
4. `configs.Get()` dipakai di mana saja untuk mengambil config yang sudah diinisialisasi

---

## 6. Database & Migrasi

### Skema Tabel

#### Tabel `users`

| Kolom        | Tipe          | Keterangan                     |
|--------------|---------------|--------------------------------|
| `id`         | BIGINT PK AI  | Primary key, auto increment    |
| `email`      | VARCHAR(255)  | Unique, tidak boleh null       |
| `password`   | VARCHAR(200)  | Disimpan sebagai bcrypt hash   |
| `username`   | VARCHAR(...)  | Ditambah via migrasi ke-2      |
| `created_at` | TIMESTAMP     | Default current timestamp      |
| `updated_at` | TIMESTAMP     | Default current timestamp      |
| `created_by` | VARCHAR(50)   | Email user pembuat             |
| `updated_by` | VARCHAR(50)   | Email user pengupdate          |

#### Tabel `posts`

| Kolom           | Tipe         | Keterangan                                        |
|-----------------|--------------|---------------------------------------------------|
| `id`            | INT PK AI    |                                                   |
| `user_id`       | INT NOT NULL | FK ke `users.id`                                  |
| `post_title`    | VARCHAR(255) |                                                   |
| `post_content`  | LONGTEXT     |                                                   |
| `post_hashtags` | LONGTEXT     | Disimpan sebagai string CSV: `"golang,backend"`   |
| `created_at`    | TIMESTAMP    |                                                   |
| `updated_at`    | TIMESTAMP    |                                                   |
| `created_by`    | VARCHAR(50)  |                                                   |
| `updated_by`    | VARCHAR(50)  |                                                   |

#### Tabel `comments`

| Kolom             | Tipe         | Keterangan           |
|-------------------|--------------|----------------------|
| `id`              | INT PK AI    |                      |
| `post_id`         | INT NOT NULL | FK ke `posts.id`     |
| `user_id`         | BIGINT       | FK ke `users.id`     |
| `comment_content` | LONGTEXT     |                      |
| `created_at`      | TIMESTAMP    |                      |
| `updated_at`      | TIMESTAMP    |                      |
| `created_by`      | VARCHAR(50)  |                      |
| `updated_by`      | VARCHAR(50)  |                      |

#### Tabel `activities`

| Kolom        | Tipe         | Keterangan                       |
|--------------|--------------|----------------------------------|
| `id`         | INT PK AI    |                                  |
| `post_id`    | INT NOT NULL | FK ke `posts.id`                 |
| `user_id`    | BIGINT       | FK ke `users.id`                 |
| `is_liked`   | BOOLEAN      | Status like/unlike dari user     |
| `created_at` | TIMESTAMP    |                                  |
| `updated_at` | TIMESTAMP    |                                  |
| `created_by` | VARCHAR(50)  |                                  |
| `updated_by` | VARCHAR(50)  |                                  |

#### Tabel `refresh_tokens`

| Kolom           | Tipe         | Keterangan                         |
|-----------------|--------------|------------------------------------|
| `id`            | INT PK AI    |                                    |
| `user_id`       | BIGINT       | FK ke `users.id`                   |
| `refresh_token` | TEXT         | Token random 36 karakter hex       |
| `expired_at`    | TIMESTAMP    | Waktu kadaluarsa (7 hari)          |
| `created_at`    | TIMESTAMP    |                                    |
| `updated_at`    | TIMESTAMP    |                                    |
| `created_by`    | VARCHAR(50)  |                                    |
| `updated_by`    | VARCHAR(50)  |                                    |

### Relasi Antar Tabel

```
users (1) ──< posts (1) ──< comments
                    └──────< activities
users (1) ──< refresh_tokens
users (1) ──< comments
users (1) ──< activities
```

### Perintah Migrasi (Makefile)

```bash
# Membuat file migrasi baru
make migrate-create name=nama_migrasi

# Jalankan semua migrasi ke atas (up)
make migrate-up

# Rollback semua migrasi (down)
make migrate-down
```

> **Catatan:** Makefile menggunakan URL `mysql://root:secret@tcp(localhost:3306)/db-simple-blog`. Sesuaikan untuk environment masing-masing.

---

## 7. Alur Inisialisasi (main.go)

```
main()
  │
  ├── 1. gin.Default()
  │         → Inisialisasi Gin router
  │
  ├── 2. configs.Init()
  │         → Baca config.yaml via Viper
  │         └── configs.Get() → *Config
  │
  ├── 3. internalsql.Connect(dsn)
  │         → Buka koneksi MySQL (*sql.DB)
  │
  ├── 4. internalsql.RunMigration(db, "./scripts/migrations")
  │         → Jalankan semua migration yang belum diaplikasikan
  │         → Jika ErrNoChange → log info (bukan fatal)
  │         → Jika error lain  → log.Fatal (app berhenti)
  │
  ├── 5. Pasang Middleware Global
  │         ├── middleware.CorsMiddleware()   → Izinkan CORS dari localhost:3000
  │         ├── gin.Logger()                  → Log setiap request
  │         └── gin.Recovery()                → Recover dari panic
  │
  ├── 6. Inisialisasi Repository
  │         ├── membershipsRepo.NewMembershipsRepository(db)
  │         └── postsRepo.NewPostsRepository(db)
  │
  ├── 7. Inisialisasi Service
  │         ├── membershipsSrv.NewMembershipsService(cfg, membershipsRepo)
  │         └── postsSrv.NewPostsService(cfg, postsRepo)
  │
  ├── 8. Inisialisasi Handler & Registrasi Route
  │         ├── membershipsHandler.RegisterRoute()  → /memberships/*
  │         └── postsHandler.RegisterRoute()         → /posts/* (dengan auth middleware)
  │
  └── 9. r.Run(cfg.Service.Port)
            → Server berjalan di port yang dikonfigurasi (default :9888)
```

---

## 8. Middleware

### `CorsMiddleware()` — `internal/middleware/cors.go`

Diaplikasikan secara **global** di semua route.

| Header                             | Nilai                                                         |
|------------------------------------|---------------------------------------------------------------|
| `Access-Control-Allow-Origin`      | `http://localhost:3000`                                       |
| `Access-Control-Allow-Credentials` | `true`                                                        |
| `Access-Control-Allow-Methods`     | `POST, OPTIONS, GET, PUT, DELETE`                             |
| `Access-Control-Allow-Headers`     | Content-Type, Authorization, Cookie, dst.                     |
| Preflight `OPTIONS`                | Dibalas dengan status `204 No Content`                        |

> ⚠️ **Perhatian:** Origin saat ini hardcoded ke `http://localhost:3000`. Harus diubah untuk production.

---

### `AuthMiddleware()` — `internal/middleware/middleware.go`

Memvalidasi JWT dari **HTTP Header `Authorization`**.

```
1. Ambil nilai header Authorization
2. Jika kosong → 401 Unauthorized + ErrMissingToken
3. jwt.ValidateToken(header, secretKey)
4. Jika gagal  → 401 Unauthorized + ErrInvalidToken
5. Jika sukses → set "id" dan "username" ke Gin context
```

> **Catatan:** `AuthMiddleware()` **tidak digunakan** pada route manapun saat ini. Route aktif menggunakan `AuthMiddlewareToken()`.

---

### `AuthMiddlewareToken()` — `internal/middleware/middleware.go`

Memvalidasi JWT dari **HTTP Cookie `access_token`**. Digunakan pada seluruh route grup `/posts`.

```
1. Ambil cookie "access_token" dari request
2. Jika tidak ada / kosong → 401 Unauthorized + ErrMissingToken
3. jwt.ValidateToken(accessToken, secretKey)
4. Jika gagal  → 401 Unauthorized + ErrInvalidToken
5. Jika sukses → set "id" dan "username" ke Gin context
```

---

## 9. Fitur: Memberships (Autentikasi)

### Route Group: `/memberships` (Public — Tidak perlu autentikasi)

---

### `POST /memberships/sign-up` — Registrasi User

**File:** `internal/handlers/memberships/login.go` → `SignUp`

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "secret123"
}
```

**Alur:**
```
Handler SignUp
  │
  ├── Bind JSON → SignUpRequest
  │
  └── Service.SignUp(ctx, request)
        ├── Repo.GetUser(email, username)
        │       └── Jika sudah ada → return ErrUsernameOrEmailAlreadyExists
        │
        ├── bcrypt.GenerateFromPassword(password)
        │       → hash password dengan cost default
        │
        └── Repo.CreateUser(model)
                → INSERT INTO users (id, email, password, username, ...)
```

**HTTP Response:**

| Status | Kondisi                        | Body                              |
|--------|--------------------------------|-----------------------------------|
| `201`  | Sukses                         | `{ "Message": "success created data" }` |
| `400`  | Request tidak valid            | `{ "message": "..." }`            |
| `409`  | Email/username sudah terdaftar | `{ "message": "Username or Email Already Exists!" }` |
| `500`  | Internal server error          | `{ "message": "..." }`            |

---

### `POST /memberships/sign-in` — Login User

**File:** `internal/handlers/memberships/login.go` → `SignIn`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

**Alur:**
```
Handler SignIn
  │
  ├── Bind JSON → SignInRequest
  │
  └── Service.SignIn(ctx, request)
        ├── Repo.GetUserByEmail(email)
        │       └── Jika nil → return ErrDataNotFound
        │
        ├── bcrypt.CompareHashAndPassword(hash, password)
        │       └── Jika salah → return ErrInvalidPassword
        │
        ├── jwt.CreateToken(id, username, secretKey)
        │       → Access token JWT (expires: 1 menit)
        │
        ├── Repo.GetRefreshToken(userId, now)
        │       └── Jika masih valid → pakai refresh token yang ada di DB
        │
        └── Jika tidak ada refresh token valid:
              ├── token.GenerateRefreshToken() → random hex 36 char
              └── Repo.InsertRefreshToken(model)
                    → INSERT INTO refresh_tokens
                      expired_at = now + 7 hari
  │
  ├── utils.SetAccessTokenCookie(c, token)
  │       → Set cookie "access_token" (MaxAge: 1 menit)
  │
  └── Response:
        {
          "message": "Success Login",
          "access_token": "<jwt>",
          "refresh_token": "<random_hex>"
        }
```

**HTTP Response:**

| Status | Kondisi                    |
|--------|----------------------------|
| `200`  | Sukses login               |
| `400`  | Password salah             |
| `404`  | Email tidak ditemukan      |
| `500`  | Internal server error      |

---

### `POST /memberships/refresh` — Refresh Access Token

**File:** `internal/handlers/memberships/refresh_token.go` → `Refresh`

**Request Body:**
```json
{
  "token": "refresh_token_string_disini"
}
```

**Alur:**
```
Handler Refresh
  │
  ├── Bind JSON → RefreshTokenRequest
  │
  ├── Service.GetIdRefreshToken(request)
  │       └── Repo.GetIdRefreshToken(token)
  │               → SELECT id, user_id, expired_at FROM refresh_tokens WHERE refresh_token = ?
  │
  └── Service.ValidateRefreshToken(userId, request)
        ├── Repo.GetRefreshToken(userId, now)
        │       → Cek expired_at >= now (belum kadaluarsa)
        │
        ├── Bandingkan request.Token == DB.RefreshToken
        │       → Jika tidak cocok → return error "invalid refresh token"
        │
        ├── Repo.GetUserById(userId)
        │       → Ambil data user untuk isi claims JWT baru
        │
        └── jwt.CreateToken(id, username, secretKey)
                → Access token baru
```

**HTTP Response:**

| Status | Kondisi                            |
|--------|------------------------------------|
| `200`  | `{ "access_token": "<jwt baru>" }` |
| `401`  | Token expired atau invalid         |
| `500`  | Internal server error              |

---

### `GET /memberships/get-user` — Ambil Data User

**File:** `internal/handlers/memberships/get_user.go` → `GetUser`

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe"
}
```

Mencari user berdasarkan **email OR username**.

**HTTP Response:**

| Status | Kondisi               |
|--------|-----------------------|
| `200`  | Data user ditemukan   |
| `404`  | User tidak ditemukan  |
| `500`  | Internal server error |

---

## 10. Fitur: Posts (Blog)

### Route Group: `/posts`

> 🔒 **Semua route dilindungi oleh `AuthMiddlewareToken()`** — membutuhkan cookie `access_token` yang valid.

---

### `POST /posts/create-post` — Buat Post Baru

**File:** `internal/handlers/posts/post.go` → `CreatePost`

**Request Body:**
```json
{
  "post_title": "Judul Artikel",
  "post_content": "Isi konten artikel...",
  "post_hashtags": ["golang", "backend", "api"]
}
```

**Alur:**
```
Handler CreatePost
  │
  ├── Bind JSON → PostRequest
  ├── Ambil userId dari Gin context (diisi oleh AuthMiddlewareToken)
  │
  └── Service.CreatePost(ctx, userId, request)
        ├── Join hashtags: ["golang","backend"] → "golang,backend" (CSV)
        ├── Buat PostModel dengan timestamps (time.Now())
        └── Repo.CreatePost(model)
                → INSERT INTO posts (user_id, post_title, post_content, post_hashtags, ...)
```

**HTTP Response:**

| Status | Kondisi                |
|--------|------------------------|
| `201`  | `{ "message": "Success Create Post" }` |
| `400`  | Request tidak valid    |
| `500`  | Internal server error  |

---

### `GET /posts/get-all-post` — Ambil Semua Post (Paginated)

**File:** `internal/handlers/posts/post.go` → `GetAllPost`

**Query Parameters:**

| Parameter   | Tipe | Keterangan                          |
|-------------|------|-------------------------------------|
| `pageIndex` | int  | Halaman ke-berapa (dimulai dari 1)  |
| `pageSize`  | int  | Jumlah item per halaman             |

**Contoh:** `GET /posts/get-all-post?pageIndex=1&pageSize=10`

**Alur:**
```
Handler GetAllPost
  │
  ├── Parse query: pageIndex, pageSize
  │
  └── Service.GetAllPost(ctx, pageSize, pageIndex)
        ├── limit  = pageSize
        ├── offset = pageSize * (pageIndex - 1)
        └── Repo.GetAllPost(limit, offset)
                → SELECT p.*, u.username FROM posts
                  JOIN users ON p.user_id = u.id
                  LIMIT ? OFFSET ?
                  (hashtags di-split dari CSV ke []string)
```

**Response Body:**
```json
{
  "message": {
    "data": [
      {
        "id": 1,
        "user_id": 1,
        "username": "johndoe",
        "post_title": "Judul Artikel",
        "post_content": "Isi...",
        "post_hashtags": ["golang", "backend"],
        "is_liked": false,
        "created_at": "2026-07-23 10:00:00",
        "updated_at": "2026-07-23 10:00:00"
      }
    ],
    "pagination": {
      "limit": 10,
      "offset": 0
    }
  }
}
```

---

### `GET /posts/get-post-by-id/:postId` — Ambil Detail Post

**File:** `internal/handlers/posts/post.go` → `GetPostById`

**URL Param:** `postId` (integer)

**Alur:**
```
Handler GetPostById
  │
  ├── Parse postId dari URL param
  │
  └── Service.GetPostById(ctx, postId)
        │
        ├── Repo.GetPostById(id)
        │       → SELECT p.*, u.username, act.is_liked FROM posts
        │         JOIN users ON p.user_id = u.id
        │         JOIN activities ON p.id = act.post_id
        │         WHERE p.id = ?
        │
        ├── Repo.CountLikedByPostID(id)
        │       → SELECT COUNT(id) FROM activities
        │         WHERE post_id = ? AND is_liked = true
        │
        └── Repo.GetCommentById(id)
                → SELECT c.*, u.username FROM comments
                  JOIN users ON c.user_id = u.id
                  WHERE c.post_id = ?
```

**Response Body:**
```json
{
  "data": {
    "detail_post": {
      "id": 1,
      "user_id": 1,
      "username": "johndoe",
      "post_title": "Judul",
      "post_content": "Isi...",
      "post_hashtags": ["golang"],
      "is_liked": true,
      "created_at": "2026-07-23 10:00:00",
      "updated_at": "2026-07-23 10:00:00"
    },
    "liked_count": 5,
    "comments": [
      {
        "id": 1,
        "user_id": 2,
        "username": "jane",
        "comment_content": "Artikel yang bagus!"
      }
    ]
  }
}
```

**HTTP Response:**

| Status | Kondisi               |
|--------|-----------------------|
| `200`  | Data ditemukan        |
| `400`  | postId tidak valid    |
| `404`  | Post tidak ditemukan  |
| `500`  | Internal server error |

---

### `POST /posts/create-comment/:postId` — Buat Komentar

**File:** `internal/handlers/posts/comment.go` → `CreateComment`

**URL Param:** `postId` (integer)

**Request Body:**
```json
{
  "comment_content": "Komentar saya di sini..."
}
```

**Alur:**
```
Handler CreateComment
  │
  ├── Bind JSON → CommentRequest
  ├── Parse postId dari URL param
  ├── Ambil userId dari Gin context (middleware)
  │
  └── Service.CreateComment(ctx, postId, userId, request)
        ├── Buat CommentModel dengan timestamps
        └── Repo.CreateComment(model)
                → INSERT INTO comments (post_id, user_id, comment_content, ...)
```

**HTTP Response:**

| Status | Kondisi                           |
|--------|-----------------------------------|
| `201`  | `{ "message": "success create comment" }` |
| `400`  | postId tidak valid / request error|
| `401`  | Tidak terautentikasi              |
| `404`  | Post tidak ditemukan              |
| `500`  | Internal server error             |

---

### `POST /posts/user-activity/:postId` — Like / Unlike Post

**File:** `internal/handlers/posts/activities.go` → `InsertUpdateActivities`

**URL Param:** `postId` (integer)

**Request Body:**
```json
{
  "is_liked": true
}
```

**Alur:**
```
Handler InsertUpdateActivities
  │
  ├── Bind JSON → ActivityRequest
  ├── Parse postId dari URL param
  ├── Ambil userId dari Gin context (middleware)
  │
  └── Service.InsertUpdateActivities(ctx, postId, userId, request)
        │
        ├── Repo.GetActivities(postId, userId)
        │       → Cek apakah sudah ada record untuk user + post ini
        │
        ├── Jika TIDAK ada record
        │       └── Repo.CreateActivities(model)
        │               → INSERT INTO activities (post_id, user_id, is_liked, ...)
        │
        └── Jika ADA record
                └── Repo.UpdateActivities(model)
                        → UPDATE activities SET is_liked = ?, updated_at = ?, updated_by = ?
                          WHERE post_id = ? AND user_id = ?
```

**HTTP Response:**

| Status | Kondisi               |
|--------|-----------------------|
| `200`  | `{ "message": "success" }` |
| `400`  | postId tidak valid    |
| `401`  | Tidak terautentikasi  |
| `500`  | Internal server error |

---

## 11. Package Shared (pkg/)

### `pkg/internalsql` — Koneksi Database & Migration

```go
// Membuka koneksi MySQL, fatal jika gagal
func Connect(datasource string) (*sql.DB, error)

// Menjalankan semua file migration yang belum diaplikasikan secara otomatis.
// migrationPath: path relatif ke folder migration (contoh: "./scripts/migrations")
// - Menggunakan golang-migrate dengan source driver "file" dan database driver "mysql"
// - Jika ErrNoChange: log info (tidak fatal, berarti semua sudah up-to-date)
// - Jika error lain: log.Fatal (aplikasi berhenti)
func RunMigration(db *sql.DB, migrationPath string)
```

**Dependency tambahan yang digunakan:**

| Package                                          | Fungsi                                      |
|--------------------------------------------------|---------------------------------------------|
| `github.com/golang-migrate/migrate/v4`           | Core library golang-migrate                 |
| `github.com/golang-migrate/migrate/v4/database/mysql` | Driver database MySQL untuk migrate    |
| `github.com/golang-migrate/migrate/v4/source/file`   | Driver source untuk membaca file `.sql` |

---

### `pkg/jwt` — JSON Web Token

```go
// Membuat JWT access token dengan algoritma HS256
// Claims: id, username, exp (1 menit dari sekarang)
func CreateToken(id int, username, secretKey string) (string, error)

// Memvalidasi token dan mengekstrak claims
func ValidateToken(tokenStr, secretKey string) (id int, username string, err error)
```

**Detail JWT:**
- **Algoritma:** HS256 (HMAC SHA-256)
- **Expiry Access Token:** 1 menit
- **Claims:** `{ "id": int, "username": string, "exp": unix_timestamp }`

---

### `pkg/token` — Refresh Token Generator

```go
// Generate string random 36 karakter menggunakan crypto/rand
func GenerateRefreshToken() string
```

Menggunakan `crypto/rand` (cryptographically secure random), menghasilkan 18 bytes yang di-hex-encode menjadi 36 karakter.

---

### `pkg/utils` — Utility Helpers

#### `JsonTime`
Custom type untuk format timestamp di JSON. Semua waktu diserialisasi dalam format `"2006-01-02 15:04:05"` (bukan ISO 8601 standar).

```go
type JsonTime time.Time
// Format: "2006-01-02 15:04:05"
```

#### `SetAccessTokenCookie`
```go
// Set cookie "access_token" dengan MaxAge 1 menit
func SetAccessTokenCookie(c *gin.Context, token string)
```

#### `SetCookie`
```go
// Fungsi generik untuk set cookie
// HttpOnly: false (dapat diakses oleh JavaScript)
func SetCookie(c *gin.Context, name, token string, maxExpired time.Duration)
```

---

## 12. Error Constants

Seluruh error sentral didefinisikan di `internal/constants/errors.go`:

| Konstanta                           | Pesan Error                           | Konteks Penggunaan                    |
|-------------------------------------|---------------------------------------|---------------------------------------|
| `ErrUsernameOrEmailAlreadyExists`   | `"Username or Email Already Exists!"` | Sign up duplikat                      |
| `ErrMissingToken`                   | `"Missing Token"`                     | Request tanpa token auth              |
| `ErrUnauthorized`                   | `"Unauthorized"`                      | Akses resource tanpa izin             |
| `ErrDataNotFound`                   | `"Data Not Found"`                    | User tidak ditemukan saat sign in     |
| `ErrPostNotFound`                   | `"Error Post Not Found"`              | Post tidak ditemukan                  |
| `ErrInvalidToken`                   | `"invalid token"`                     | JWT tidak valid atau expired          |
| `ErrInvalidPassword`                | `"invalid password"`                  | Password salah saat sign in           |
| `ErrTokenExpired`                   | `"refresh token has expired"`         | Refresh token sudah kadaluarsa        |

---

## 13. Alur Request End-to-End

### Skenario: User Login → Buat Post → Token Expired → Refresh

```
┌─────────────────────────────────────────────────────────────┐
│ STEP 1 — Registrasi (jika belum punya akun)                 │
│                                                             │
│  POST /memberships/sign-up                                  │
│  { email, username, password }                              │
│                ↓                                            │
│  201 Created                                                │
└─────────────────────────────────────────────────────────────┘
                ↓
┌─────────────────────────────────────────────────────────────┐
│ STEP 2 — Login                                              │
│                                                             │
│  POST /memberships/sign-in                                  │
│  { email, password }                                        │
│                ↓                                            │
│  200 OK + Set Cookie "access_token" (1 menit)              │
│  Response: { access_token, refresh_token }                  │
└─────────────────────────────────────────────────────────────┘
                ↓
┌─────────────────────────────────────────────────────────────┐
│ STEP 3 — Akses Protected Route                              │
│  (Browser otomatis kirim cookie "access_token")             │
│                                                             │
│  POST /posts/create-post                                    │
│  { post_title, post_content, post_hashtags }                │
│                ↓                                            │
│  AuthMiddlewareToken → validasi cookie → inject userId      │
│                ↓                                            │
│  201 Created                                                │
└─────────────────────────────────────────────────────────────┘
                ↓
┌─────────────────────────────────────────────────────────────┐
│ STEP 4 — Token Expired (setelah 1 menit)                    │
│                                                             │
│  POST /memberships/refresh                                  │
│  { token: "<refresh_token>" }                               │
│                ↓                                            │
│  Validasi refresh token di DB (expired_at >= now)           │
│  Generate JWT baru                                          │
│                ↓                                            │
│  200 OK: { access_token: "<jwt baru>" }                     │
│  (Client harus update cookie secara manual)                 │
└─────────────────────────────────────────────────────────────┘
```

---

## 14. Ringkasan Endpoint API

### 🔓 Public Routes

| Method | Endpoint                   | Fungsi                              |
|--------|----------------------------|-------------------------------------|
| POST   | `/memberships/sign-up`     | Registrasi user baru                |
| POST   | `/memberships/sign-in`     | Login, mendapat access + refresh token |
| GET    | `/memberships/get-user`    | Ambil data user by email/username   |
| POST   | `/memberships/refresh`     | Refresh access token                |

### 🔒 Protected Routes (Cookie `access_token` wajib)

| Method | Endpoint                              | Fungsi                                          |
|--------|---------------------------------------|-------------------------------------------------|
| POST   | `/posts/create-post`                  | Buat post baru                                  |
| GET    | `/posts/get-all-post`                 | Daftar semua post (paginated)                   |
| GET    | `/posts/get-post-by-id/:postId`       | Detail post + jumlah like + komentar            |
| POST   | `/posts/create-comment/:postId`       | Buat komentar pada post                         |
| POST   | `/posts/user-activity/:postId`        | Like atau unlike sebuah post                    |

---

## 15. Catatan Penting & Perhatian

### ⚠️ Hal yang Perlu Diperhatikan

1. **JWT Access Token expire sangat cepat (1 menit)**  
   Nilai ini sangat ketat. Untuk production pertimbangkan 15 menit – 1 jam.

2. **Cookie `HttpOnly: false`**  
   Cookie `access_token` bisa dibaca oleh JavaScript → rentan XSS.  
   Pertimbangkan mengubah ke `HttpOnly: true`.

3. **CORS Origin hardcoded**  
   `Access-Control-Allow-Origin` hardcoded ke `http://localhost:3000`.  
   Harus dikonfigurasi dinamis via environment variable untuk staging/production.

4. **`GetPostById` membutuhkan record di tabel `activities`**  
   Query JOIN ke `activities` menggunakan INNER JOIN, sehingga post yang belum pernah di-like/unlike oleh siapapun akan mengembalikan `ErrPostNotFound`. Perlu diubah ke LEFT JOIN.

5. **`AuthMiddleware()` tidak digunakan**  
   Fungsi validasi token via Header `Authorization` sudah dibuat tapi tidak dipasang ke route manapun. Yang aktif adalah `AuthMiddlewareToken()` via Cookie.

6. **Hashtag disimpan sebagai CSV string**  
   Di database `post_hashtags` tersimpan sebagai `"golang,backend,api"`. Konversi ke/dari `[]string` dilakukan di layer repository menggunakan `strings.Join` dan `strings.Split`.

7. **`GetUser` handler menggunakan struct `SignUpRequest`**  
   Endpoint `GET /memberships/get-user` menggunakan `SignUpRequest` (yang memiliki field `password`) padahal password tidak dibutuhkan. Ini inkonsistensi naming yang perlu dirapikan.

---

## 16. Changelog Dokumentasi

| Versi | Tanggal    | Perubahan                                                                 | Oleh             |
|-------|------------|---------------------------------------------------------------------------|------------------|
| 1.0.0 | 2026-07-23 | Dokumentasi awal — analisis seluruh codebase                              | Antigravity (AI) |
| 1.1.0 | 2026-07-24 | Tambah fitur auto-migration: `RunMigration` di `pkg/internalsql`, dipanggil di `main.go` saat startup | Antigravity (AI) |
