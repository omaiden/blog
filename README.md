# 📝 Blog Service

A backend service for managing blog posts, built with Go and PostgreSQL.

---

## 🚀 Getting Started

### 📦 Prerequisites

- Go 1.24
- PostgreSQL running at `localhost:5432`
- Database `blog` must exist

---

## Docker

- To use postgres in docker

```sh
docker-compose up
```

## ⚙️ Environment Variables

Sets the following environment variables by default:

```bash
DB_URL=postgres://user:password@localhost:5432/blog?sslmode=disable
TEST_DB_URL=postgres://user:passwordw@localhost:5432/%s?sslmode=disable
```

## Run test in local

```sh
make test
```

## Run the service

```sh
make run
```
