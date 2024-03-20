CREATE TABLE
    IF NOT EXISTS post_comments (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        postId UUID NOT NULL,
        userId UUID NOT NULL,
        comment TEXT NOT NULL,
        createdAt TIMESTAMP DEFAULT now (),
        updatedAt TIMESTAMP DEFAULT now (),

        CONSTRAINT fk_post FOREIGN KEY (postId) REFERENCES posts (id) ON DELETE CASCADE,
        CONSTRAINT fk_user FOREIGN KEY (userId) REFERENCES users (id) ON DELETE CASCADE
    );

CREATE TRIGGER update_post_comments_updated_at
    BEFORE UPDATE
    ON post_comments
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated();