CREATE TABLE chats (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type          integer NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE chat_members (
    id            integer PRIMARY KEY AUTOINCREMENT,
    chat_id       UUID NOT NULL REFERENCES chats(id),
    user_id       UUID NOT NULL REFERENCES users(id),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);