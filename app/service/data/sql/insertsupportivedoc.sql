INSERT INTO supportive_docs (
    board_id, card_id, name, title, category,
    labels, attachments, comments, updated_at   
) VALUES (
    :board_id, :card_id, :name, :title, :category,
    :labels, :attachments, :comments, NOW()
)
RETURNING id
