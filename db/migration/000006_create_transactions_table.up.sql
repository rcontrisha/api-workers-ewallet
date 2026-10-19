CREATE TABLE transaction (
    id                  uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    trx_code            varchar     NOT NULL UNIQUE,
    type                varchar     NOT NULL,
    status              varchar     NOT NULL,
    amount              bigint      NOT NULL,
    tax                 bigint      NOT NULL DEFAULT 0,
    discount            bigint      NOT NULL DEFAULT 0,
    total               bigint      NOT NULL,
    sender_wallet_id    uuid        REFERENCES wallet (id),
    receiver_wallet_id  uuid        NOT NULL REFERENCES wallet (id),
    payment_method_id   uuid        REFERENCES payment_method (id),
    note                varchar,
    deleted_by_sender   boolean     NOT NULL DEFAULT false,
    deleted_by_receiver boolean     NOT NULL DEFAULT false,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz
);

CREATE INDEX idx_transaction_sender_created   ON transaction (sender_wallet_id, created_at);
CREATE INDEX idx_transaction_receiver_created ON transaction (receiver_wallet_id, created_at);