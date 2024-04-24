# BinaryBlog - Microservice edition
  
![Go Version](https://img.shields.io/github/go-mod/go-version/jay-bhogayata/blog)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Last Commit](https://img.shields.io/github/last-commit/jay-bhogayata/blog)
![Issues](https://img.shields.io/github/issues/jay-bhogayata/blog)
![Stars](https://img.shields.io/github/stars/jay-bhogayata/blog)

<br/>

<div align="center">
<img src="./images/logo.jpeg" alt="BinaryBlog" width="200" height="200"/>
</div>

## Description

This repo contains microservices of blog application.

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
# blog-service
go run cmd/api/blog/main.go [option]
# [option] = 1, 2

# 1: Database migration-up
# 2: Database migration-down

# notification-service
go run .

```

> Keep in mind that you need to have a postgres database running on your local machine or in a docker container before running the above command

> leave `[option]` empty to run the server