# Blog Backend

## Description

This is a simple blog backend that allows you to create, read, update and delete blog posts. It is built using goland and postgres.

## Requirements

- Go 1.21.5 and above
- Postgres 13.3 and above
- Docker 20.10.7 and above
- sqlc 1.25.0 and above
- golang-migrate 4.17.0 and above

## Sqlc

```bash
# Docker recommended
docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate
```
or

Visit [sqlc](https://sqlc.dev/docs/install) to install sqlc

## Usage

```bash
go run cmd/api/blog/main.go [option]
# [option] = 1, 2, 3

# 1: Start blog server
# 2: Database migration-up
# 3: Database migration-down

```