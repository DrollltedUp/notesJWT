
CREATE TABLE users (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email         TEXT        NOT NULL UNIQUE CHECK (email = lower(email)),
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE notes (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    owner_id   BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title      TEXT        NOT NULL CHECK (char_length(title) BETWEEN 1 AND 200),
    body       TEXT        NOT NULL DEFAULT '' CHECK (char_length(body) <= 10000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX notes_owner_id_idx ON notes (owner_id, id);