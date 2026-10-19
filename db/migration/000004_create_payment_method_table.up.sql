CREATE TABLE payment_method (
    id                    uuid    PRIMARY KEY DEFAULT gen_random_uuid(),
    gateway_payment_channel varchar NOT NULL UNIQUE,
    name                  varchar NOT NULL UNIQUE,
    logo                  varchar,
    is_active             boolean NOT NULL DEFAULT true
);