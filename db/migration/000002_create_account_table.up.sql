CREATE TABLE account (
    id          uuid        PRIMARY KEY REFERENCES users (id),
    fullname    varchar     NOT NULL,
    phone       varchar     NOT NULL UNIQUE,
    avatar      varchar,
    pin         varchar,
    is_verified boolean     NOT NULL DEFAULT false,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz
);