CREATE TABLE "user"
(
    id            int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    email         varchar(50)  NOT NULL,
    first_name    varchar(50)  NULL,
    last_name     varchar(50)  NULL,
    middle_name   varchar(50)  NULL,
    password_hash varchar(512) NOT NULL,
    created_at    timestamp    NOT NULL,
    updated_at    timestamp    NOT NULL
);

create unique index table_name_user_uindex
    on "user" (email);