CREATE TABLE attachments (
    id SERIAL PRIMARY KEY,
    trello_url TEXT NOT NULL,
    storage_url TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
