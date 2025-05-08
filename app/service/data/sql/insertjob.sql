INSERT INTO jobs (
    type, state, cards, errors, started_at, completed_at
) VALUES (
    :type, :state, :cards, :errors, :started_at, :completed_at
)
RETURNING id
