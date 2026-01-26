-- +goose Up
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_notes_user
            FOREIGN KEY (user_id) REFERENCES users(id)
                   ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS notes;
