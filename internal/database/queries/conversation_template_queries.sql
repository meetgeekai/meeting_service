-- name: GetConversationTemplateNames :many
SELECT
    ct.id,
    ct.name
FROM
    conversation_templates ct
WHERE
    ct.id IN (sqlc.slice('ids'))
    AND ct.user_id = ?
;
