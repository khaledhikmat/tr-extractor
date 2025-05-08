CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    state TEXT NOT NULL,
    cards BIGINT NOT NULL,
    errors BIGINT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP
);
