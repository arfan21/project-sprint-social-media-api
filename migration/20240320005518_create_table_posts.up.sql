CREATE TABLE
    IF NOT EXISTS posts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        userId UUID NOT NULL,
        body TEXT NOT NULL,
        tags VARCHAR(255) [],
        createdAt TIMESTAMP DEFAULT now (),
        updatedAt TIMESTAMP DEFAULT now (),

        CONSTRAINT fk_user FOREIGN KEY (userId) REFERENCES users (id) ON DELETE CASCADE
    );

CREATE TRIGGER update_posts_updated_at
    BEFORE UPDATE
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated();