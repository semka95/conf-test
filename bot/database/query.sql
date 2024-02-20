-- name: CreateReport :exec
INSERT INTO report(url, title, starting_at, duration_minutes, reporters, conference_id, status)
    VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetAllReports :many
SELECT
    *
FROM
    report;

