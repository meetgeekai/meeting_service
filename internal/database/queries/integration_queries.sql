-- name: GetConnectedCalendarVendors :many
SELECT DISTINCT
    i.vendor
FROM
    integrations i
WHERE
    i.user_uuid = ?
    AND i.type = 'Calendar'
    AND i.active = 1
    AND i.error_code IS NULL
    AND i.vendor IN ('google', 'microsoft')
;
