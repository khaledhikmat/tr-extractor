INSERT INTO inheritance_confinments (
    board_id, card_id, name, title, generation,
    labels, attachments, comments, updated_at   
) VALUES (
    :board_id, :card_id, :name, :title, :generation,
    :labels, :attachments, :comments, NOW()
)
RETURNING id
