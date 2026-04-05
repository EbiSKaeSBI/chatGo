CREATE TYPE dialog_type AS ENUM ('group','private');

CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    username      TEXT                      NOT NULL UNIQUE,
    password_hash TEXT                      NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT now() NOT NULL
);


CREATE TABLE dialogs
(
    id         BIGSERIAL PRIMARY KEY,
    type       dialog_type               NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE dialog_participants
(
    user_id   BIGINT REFERENCES users (id),
    dialog_id BIGINT REFERENCES dialogs (id),
    PRIMARY KEY (user_id, dialog_id)
);

CREATE TABLE messages
(
    id         BIGSERIAL PRIMARY KEY,
    dialog_id  BIGINT REFERENCES dialogs (id) NOT NULL,
    sender_id  BIGINT REFERENCES users (id)   NOT NULL,
    text       TEXT                           NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()      NOT NULL
);

