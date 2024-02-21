-- name: CreateReport :exec
INSERT INTO report(url, title, starting_at, duration_minutes, reporters, status)
    VALUES (?, ?, ?, ?, ?, ?);

-- name: GetAllReports :many
SELECT
    starting_at,
    duration_minutes,
    title,
    reporters,
    url
FROM report
WHERE starting_at >= ?;

-- name: GetAllRatings :many
SELECT
    user.id_data,
    report.url,
    rating.content_score,
    rating.presentation_score,
    rating.notes
FROM
    rating
    LEFT JOIN report ON rating.report_id = report.id
    LEFT JOIN user ON rating.user_id = user.telegram_id
WHERE
    report.starting_at >= ?
    AND rating.rating_type = 'score'
    AND user.id_data IS NOT NULL;

-- name: AddFavoriteReport :exec
INSERT INTO favorite_reports(user_id, report_id)
    VALUES (?, ?);

-- name: RemoveFavoriteReport :exec
DELETE FROM favorite_reports
WHERE user_id = ?
    AND report_id = ?;

-- name: GetFavoriteReports :many
SELECT
    report.starting_at,
    report.duration_minutes,
    report.title,
    report.reporters,
    report.url
FROM
    favorite_reports
    LEFT JOIN report ON favorite_reports.report_id = report.report_id
    LEFT JOIN user ON favorite_reports.user_id = user.telegram_id
WHERE
    favorite_reports.user_id = ? -- TODO: filter deleted reports

-- name: UpdateUserData :exec
UPDATE
    user
SET
    id_data = ?
WHERE
    telegram_id = ?;

-- name: CreateRating :exec
INSERT INTO rating(report_id, user_id, rating_type, content_score, presentation_score, notes)
    VALUES (?, ?, ?, ?, ?, ?);

-- name: GetAllUserRatings :many
SELECT
    user.id_data,
    report.url,
    rating.content_score,
    rating.presentation_score,
    rating.notes
FROM
    rating
    LEFT JOIN report ON rating.report_id = report.id
    LEFT JOIN user ON rating.user_id = user.telegram_id
WHERE
    user.telegram_id = ?
    AND report.starting_at >= ?;

-- name: GetAllUserReportsNoScore :many
SELECT
    starting_at,
    duration_minutes,
    title,
    reporters,
    url
FROM
    report
WHERE
    report.starting_at >= ?
    AND id NOT IN (
        SELECT
            report_id
        FROM
            rating
        WHERE
            user_id = ?);

