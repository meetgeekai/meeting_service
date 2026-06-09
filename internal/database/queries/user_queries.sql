-- name: GetUserForUpcomingMeetings :one
SELECT
    u.uuid,
    u.name,
    u.email
FROM
    users u
WHERE
    u.id = ?
    AND u.deleted_at IS NULL
;
