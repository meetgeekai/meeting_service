-- name: GetUserForUpcomingMeetings :one
SELECT
    u.id,
    u.name,
    u.email
FROM
    users u
WHERE
    u.uuid = ?
    AND u.deleted_at IS NULL
;
