INSERT INTO properties (
    board_id, card_id, name, location_ar, location_en, lot, type, status, owner, area, shares,
    is_organized, is_effects, labels, attachments, comments, updated_at   
) VALUES (
    :board_id, :card_id, :name, :location_ar, :location_en, :lot, :type, :status, :owner, :area, :shares,
    :is_organized, :is_effects, :labels, :attachments, :comments, NOW()
)
RETURNING id
