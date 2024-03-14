-- name: CreateCategory :one
INSERT INTO categories (name) VALUES ($1) RETURNING *;

-- name: GetCategoryById :one
SELECT * FROM categories WHERE id = $1;

-- name: GetAllCategories :many
SELECT * FROM categories;

-- name: UpdateCategory :one
UPDATE categories SET name = $1 WHERE id = $2 RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, is_verified, verification_token) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE user_id = $1;

-- name: GetUserByVerificationToken :one
SELECT * FROM users WHERE verification_token = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetAllUsers :many
SELECT user_id,username,email,is_verified FROM users;

-- name: UpdateUser :one
UPDATE users SET username = $1, email = $2, password_hash = $3, is_verified = $4, verification_token = $5 WHERE user_id = $6 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = $1;

-- name: VerifyUser :exec
UPDATE users SET is_verified = true, verification_token = null , updated_at = now() WHERE verification_token = $1  RETURNING *;


-- name: CreateArticle :one
INSERT INTO articles (title, content, category_id, user_id) VALUES ($1, $2, $3, $4) RETURNING article_id;

-- name: PublishArticle :exec
UPDATE articles SET is_published = true, updated_at = now() WHERE article_id = $1;

-- name: GetAllArticles :many
SELECT a.article_id, a.title, a.content, a.user_id, c.name as category_name, a.created_at, a.updated_at, a.is_published
FROM articles a
LEFT JOIN categories c ON a.category_id = c.id
WHERE a.is_published = true;


-- name: GetAllArticleByUser :many
SELECT a.article_id, a.title, a.content, a.user_id, c.name as category_name, a.created_at, a.updated_at, a.is_published
FROM articles a
LEFT JOIN categories c ON a.category_id = c.id
WHERE a.user_id = $1;

-- name: GetUserIdByArticleId :one
SELECT user_id FROM articles WHERE article_id = $1;