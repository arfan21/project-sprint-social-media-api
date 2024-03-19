CREATE TABLE
    IF NOT EXISTS friends (
        userIdAdder UUID NOT NULL,
        userIdAdded UUID NOT NULL,
        createdAt TIMESTAMP DEFAULT now (),
        updatedAt TIMESTAMP DEFAULT now (),

        PRIMARY KEY (userIdAdder, userIdAdded),
        CONSTRAINT fk_userIdAdder FOREIGN KEY (userIdAdder) REFERENCES users (id) ON DELETE CASCADE,
        CONSTRAINT fk_userIdAdded FOREIGN KEY (userIdAdded) REFERENCES users (id) ON DELETE CASCADE
    );

CREATE TRIGGER update_friends_updated_at
  BEFORE UPDATE
  ON friends
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_set_updated();