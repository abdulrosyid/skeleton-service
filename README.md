# skeleton-service

Template service Golang dengan arsitektur mirip `new_sikda` (hexagonal sederhana), plus sample use case:
- User Registrasi
- User Login (JWT)

## Endpoint
- `GET /health`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/profile` (Bearer token)

## Contoh request
Register:
```bash
curl --location 'http://localhost:8080/api/v1/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
  "full_name": "Budi",
  "email": "budi@mail.com",
  "password": "rahasia123"
}'
```

Login:
```bash
curl --location 'http://localhost:8080/api/v1/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
  "email": "budi@mail.com",
  "password": "rahasia123"
}'
```

Profile (protected):
```bash
curl --location 'http://localhost:8080/api/v1/profile' \
--header 'Authorization: Bearer <JWT_TOKEN_DARI_LOGIN>'
```

## Tabel MySQL sample
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    full_name VARCHAR(120) NOT NULL,
    email VARCHAR(120) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Run
```bash
go mod tidy
go run .
```

## Run with Docker (App + MySQL)
```bash
cp .env.example .env
docker compose up -d --build
```

Check status:
```bash
docker compose ps
docker logs -f myapp
curl http://localhost:8080/health
```
