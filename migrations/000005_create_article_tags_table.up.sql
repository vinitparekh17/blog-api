CREATE TABLE IF NOT EXISTS article_tags (
    article_id uuid REFERENCES articles(article_id) ON DELETE CASCADE NOT NULL,
    tag_id INT REFERENCES tags(id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (article_id, tag_id)
);