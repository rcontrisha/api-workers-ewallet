CREATE TABLE favorite (
    user_id          uuid        NOT NULL REFERENCES users (id),
    favorite_user_id uuid        NOT NULL REFERENCES users (id),
    created_at       timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, favorite_user_id)
);