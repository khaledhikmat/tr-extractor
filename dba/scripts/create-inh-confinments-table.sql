CREATE TABLE inheritance_confinments (
    id SERIAL PRIMARY KEY,
    board_id TEXT NOT NULL,
    card_id TEXT NOT NULL,
    name TEXT NOT NULL,
    title TEXT NOT NULL,
    generation NUMERIC NOT NULL,
    labels TEXT[],
    attachments TEXT[],
    comments TEXT[],
    updated_at TIMESTAMP NOT NULL
);
