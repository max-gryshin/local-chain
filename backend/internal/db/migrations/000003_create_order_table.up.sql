create table "order" (
    id              int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    status          smallserial
        CONSTRAINT state_check CHECK (status > 0 and status <= 15),
    amount          numeric       NOT NULL,
    wallet_id       int           NOT NULL,
    description     varchar(1024) NOT NULL,
    request_reasons jsonb         NOT NULL,
    created_at      timestamp     NOT NULL,
    updated_at      timestamp     NOT NULL,
    created_by      int           NOT NULL,
    updated_by      int           NOT NULL
);