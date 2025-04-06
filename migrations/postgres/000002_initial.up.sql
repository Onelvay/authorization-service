CREATE EXTENSION IF NOT EXISTS pgcrypto;


CREATE TABLE IF NOT EXISTS billings (
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id                  VARCHAR PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    iin                 VARCHAR NOT NULL,
    correlation_id      VARCHAR NOT NULL,
    invoice_id          VARCHAR NOT NULL UNIQUE,
    amount              NUMERIC NOT NULL,
    currency            VARCHAR NOT NULL,
    terminal_id         VARCHAR NOT NULL,
    description         VARCHAR NOT NULL,
    account_id          VARCHAR NULL,
    name                VARCHAR NULL,
    email               VARCHAR NULL,
    phone               VARCHAR NOT NULL,
    back_link           VARCHAR NOT NULL,
    failure_back_link   VARCHAR NULL,
    post_link           VARCHAR NOT NULL,
    failure_post_link   VARCHAR NULL,
    language            VARCHAR NULL,
    data                VARCHAR NULL,
    card_save           BOOLEAN DEFAULT FALSE,
    source              VARCHAR DEFAULT 'epay'
);

CREATE TABLE IF NOT EXISTS cards(
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id          VARCHAR PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    card_id     VARCHAR NOT NULL,
    account_id  VARCHAR NOT NULL,
    terminal_id VARCHAR NOT NULL,
    type        VARCHAR NOT NULL,
    mask        VARCHAR NOT NULL,
    issuer      VARCHAR NOT NULL,
    is_default  BOOLEAN DEFAULT FALSE
);
