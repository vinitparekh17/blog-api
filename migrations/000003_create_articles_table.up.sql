CREATE TABLE IF NOT EXISTS articles (
    article_id  uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    user_id uuid REFERENCES users(user_id) ON DELETE CASCADE NOT NULL,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    is_published BOOLEAN DEFAULT FALSE
);