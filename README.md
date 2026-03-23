# Go Auth Backend

Boshlang'ich Go dasturlash tilida tayyorlangan oddiy **Auth API backend**.
Ushbu loyiha quyidagi imkoniyatlarni beradi:

- Register (ro'yxatdan o'tish)
- Login (tizimga kirish)
- Forgot password (reset token olish)
- Reset password (token orqali parolni tiklash)
- Profilni olish
- Profil ma'lumotlarini yangilash
- Parolni o'zgartirish

Texnologiyalar:

- Go
- Gin (HTTP framework)
- GORM
- SQLite
- JWT
- bcrypt

---

## 1) Project struktura

```txt
.
├── cmd
│   └── server
│       └── main.go
├── config
│   └── config.go
├── database
│   └── database.go
├── handlers
│   ├── auth_handler.go
│   └── profile_handler.go
├── middleware
│   └── auth_middleware.go
├── models
│   ├── password_reset_token.go
│   └── user.go
├── routes
│   └── routes.go
├── utils
│   ├── jwt.go
│   ├── password.go
│   ├── response.go
│   └── token.go
├── auth-api-docs.md
├── go.mod
└── README.MD
```

---

## 2) Ishga tushirish (Local setup)

### Talablar

- Go o'rnatilgan bo'lishi kerak (`go version`)

### Qadamlar

1. Dependencylarni o'rnatish:

```bash
go mod tidy
```

2. Serverni ishga tushirish:

```bash
go run ./cmd/server
```

3. Server default quyidagi manzilda ishlaydi:

`http://localhost:8080`

---

## 3) Environment variables

Loyiha `.env` talab qilmaydi, lekin quyidagi o'zgaruvchilarni berishingiz mumkin:

- `PORT` (default: `:8080`)
- `JWT_SECRET` (default: `dev-secret-change-me`)
- `DATABASE_PATH` (default: `auth.db`)

PowerShell misol:

```powershell
$env:PORT=":8080"
$env:JWT_SECRET="super-secret-key"
$env:DATABASE_PATH="auth.db"
go run ./cmd/server
```

---

## 4) API base URL

Base URL:

`http://localhost:8080/api`

Protected endpointlar uchun header:

`Authorization: Bearer <JWT_TOKEN>`

---

## 5) Endpointlar ro'yxati

### Auth

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/forgot-password`
- `POST /auth/reset-password`

### Profile

- `GET /profile/me` (protected)
- `PUT /profile/me` (protected)
- `PUT /profile/change-password` (protected)

Batafsil request/response namunalari: `auth-api-docs.md`.

---

## 6) Tez test qilish (cURL misollar)

### Register

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"full_name\":\"Ali Valiyev\",\"email\":\"ali@mail.com\",\"phone\":\"+998901234567\",\"password\":\"123456\"}"
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"ali@mail.com\",\"password\":\"123456\"}"
```

### Forgot Password

```bash
curl -X POST http://localhost:8080/api/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"ali@mail.com\"}"
```

### Reset Password

```bash
curl -X POST http://localhost:8080/api/auth/reset-password \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"<reset_token>\",\"new_password\":\"newStrongPass123\"}"
```

### Get Profile

```bash
curl -X GET http://localhost:8080/api/profile/me \
  -H "Authorization: Bearer <jwt_token>"
```

### Update Profile

```bash
curl -X PUT http://localhost:8080/api/profile/me \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d "{\"full_name\":\"Ali Updated\",\"phone\":\"+998909999999\"}"
```

### Change Password

```bash
curl -X PUT http://localhost:8080/api/profile/change-password \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d "{\"current_password\":\"123456\",\"new_password\":\"newStrongPass123\"}"
```

---

## 7) Ma'lumotlar bazasi

- Database sifatida `SQLite` ishlatiladi.
- Fayl nomi default: `auth.db`.
- `AutoMigrate` orqali quyidagi jadvallar avtomatik yaratiladi:
  - `users`
  - `password_reset_tokens`
