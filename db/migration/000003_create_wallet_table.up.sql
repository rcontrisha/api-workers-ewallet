CREATE TABLE wallet (
    id         uuid        PRIMARY KEY REFERENCES users (id),
    balance    bigint      NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz
);