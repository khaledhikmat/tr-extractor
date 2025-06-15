UPDATE inheritance_confinments SET
    board_id = $1,
    card_id = $2,
    name = $3,
    title = $4,
    generation = $5,
    labels = $6,
    attachments = $7,
    comments = $8,
    updated_at = NOW()
WHERE id = $9;