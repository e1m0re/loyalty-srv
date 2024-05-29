-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE users
(
    id         SERIAL
        CONSTRAINT users_pk PRIMARY KEY,
    password   VARCHAR(255) NOT NULL,
    last_login TIMESTAMP,
    username   VARCHAR(50)  NOT NULL UNIQUE
);

CREATE TABLE orders
(
    id         SERIAL
        CONSTRAINT orders_pk
            PRIMARY KEY,
    "user"     INTEGER
        CONSTRAINT orders_users_id_fk
            REFERENCES users,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    status     VARCHAR(50)                         NOT NULL,
    number     VARCHAR(50)
        UNIQUE,
    accrual    INTEGER,
    calculated BOOLEAN   DEFAULT false             NOT NULL
);

CREATE TABLE accounts
(
    id      SERIAL
        CONSTRAINT accounts_pk
            PRIMARY KEY,
    "user"  INTEGER NOT NULL
        CONSTRAINT accounts_users_id_fk
            REFERENCES users,
    balance money   NOT NULL
);

CREATE TABLE accounts_changes
(
    id      SERIAL
        CONSTRAINT accounts_changes_pk
            PRIMARY KEY,
    account INTEGER   NOT NULL
        CONSTRAINT accounts_changes_accounts_id_fk
            REFERENCES accounts,
    amount  money     NOT NULL,
    ts      TIMESTAMP NOT NULL,
    "order" VARCHAR(50)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
