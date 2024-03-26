CREATE INDEX IF NOT EXISTS idx_tags ON posts USING GIN (tags);

CREATE INDEX IF NOT EXISTS idx_posts ON posts (body);