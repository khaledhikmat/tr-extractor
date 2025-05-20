INSERT INTO attachments (
    trello_url, storage_url, updated_at
) VALUES (
    :trello_url, :storage_url, :updated_at
)
RETURNING id
