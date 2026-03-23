# Auth API Docs

## Base URL

`http://localhost:8080/api`

## Authentication

Protected endpointlar uchun header:

`Authorization: Bearer <JWT_TOKEN>`

---

## 1) Register

- **Method:** `POST`
- **URL:** `/auth/register`
- **Body:**

```json
{
  "full_name": "Ali Valiyev",
  "email": "ali@mail.com",
  "phone": "+998901234567",
  "password": "123456"
}
```

- **Success (201):**

```json
{
  "message": "register successful",
  "token": "<jwt_token>",
  "user": {
    "id": 1,
    "full_name": "Ali Valiyev",
    "email": "ali@mail.com",
    "phone": "+998901234567",
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:00:00Z"
  }
}
```

---

## 2) Login

- **Method:** `POST`
- **URL:** `/auth/login`
- **Body:**

```json
{
  "email": "ali@mail.com",
  "password": "123456"
}
```

- **Success (200):**

```json
{
  "message": "login successful",
  "token": "<jwt_token>",
  "user": {
    "id": 1,
    "full_name": "Ali Valiyev",
    "email": "ali@mail.com",
    "phone": "+998901234567",
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:00:00Z"
  }
}
```

---

## 3) Forgot Password

- **Method:** `POST`
- **URL:** `/auth/forgot-password`
- **Body:**

```json
{
  "email": "ali@mail.com"
}
```

- **Success (200):**

```json
{
  "message": "if account exists, reset token generated",
  "reset_token": "<reset_token>",
  "note": "in production this token should be sent by email"
}
```

> Eslatma: Hozircha demo uchun `reset_token` response’da qaytadi.

---

## 4) Reset Password

- **Method:** `POST`
- **URL:** `/auth/reset-password`
- **Body:**

```json
{
  "token": "<reset_token>",
  "new_password": "newStrongPass123"
}
```

- **Success (200):**

```json
{
  "message": "password reset successful"
}
```

---

## 5) Get My Profile

- **Method:** `GET`
- **URL:** `/profile/me`
- **Header:** `Authorization: Bearer <jwt_token>`

- **Success (200):**

```json
{
  "user": {
    "id": 1,
    "full_name": "Ali Valiyev",
    "email": "ali@mail.com",
    "phone": "+998901234567",
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:00:00Z"
  }
}
```

---

## 6) Update My Profile

- **Method:** `PUT`
- **URL:** `/profile/me`
- **Header:** `Authorization: Bearer <jwt_token>`
- **Body:**

```json
{
  "full_name": "Ali Valiyev Updated",
  "phone": "+998909999999"
}
```

- **Success (200):**

```json
{
  "message": "profile updated",
  "user": {
    "id": 1,
    "full_name": "Ali Valiyev Updated",
    "email": "ali@mail.com",
    "phone": "+998909999999",
    "created_at": "2026-03-23T10:00:00Z",
    "updated_at": "2026-03-23T10:05:00Z"
  }
}
```

---

## 7) Change Password

- **Method:** `PUT`
- **URL:** `/profile/change-password`
- **Header:** `Authorization: Bearer <jwt_token>`
- **Body:**

```json
{
  "current_password": "123456",
  "new_password": "newStrongPass123"
}
```

- **Success (200):**

```json
{
  "message": "password changed successfully"
}
```

---

## Run Locally

1. `go mod tidy`
2. `go run ./cmd/server`

Server default: `http://localhost:8080`

## Optional Environment Variables

- `PORT` (default `:8080`)
- `JWT_SECRET` (default `dev-secret-change-me`)
- `DATABASE_PATH` (default `auth.db`)
