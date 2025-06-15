CREATE TABLE supportive_docs (
    id SERIAL PRIMARY KEY,
    board_id TEXT NOT NULL,
    card_id TEXT NOT NULL,
    name TEXT NOT NULL,
    title TEXT NOT NULL,
    category TEXT NOT NULL,
    labels TEXT[],
    attachments TEXT[],
    comments TEXT[],
    updated_at TIMESTAMP NOT NULL
);
