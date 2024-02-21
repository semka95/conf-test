// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const addFavoriteReport = `-- name: AddFavoriteReport :exec
INSERT INTO favorite_reports(user_id, report_id)
    VALUES (?, ?)
`

type AddFavoriteReportParams struct {
	UserID   int64
	ReportID int64
}

func (q *Queries) AddFavoriteReport(ctx context.Context, arg AddFavoriteReportParams) error {
	_, err := q.exec(ctx, q.addFavoriteReportStmt, addFavoriteReport, arg.UserID, arg.ReportID)
	return err
}

const createRating = `-- name: CreateRating :exec
INSERT INTO rating(report_id, user_id, rating_type, content_score, presentation_score, notes)
    VALUES (?, ?, ?, ?, ?, ?)
`

type CreateRatingParams struct {
	ReportID          int64
	UserID            int64
	RatingType        string
	ContentScore      sql.NullInt64
	PresentationScore sql.NullInt64
	Notes             sql.NullString
}

func (q *Queries) CreateRating(ctx context.Context, arg CreateRatingParams) error {
	_, err := q.exec(ctx, q.createRatingStmt, createRating,
		arg.ReportID,
		arg.UserID,
		arg.RatingType,
		arg.ContentScore,
		arg.PresentationScore,
		arg.Notes,
	)
	return err
}

const createReport = `-- name: CreateReport :exec
INSERT INTO report(url, title, starting_at, duration_minutes, reporters, conference_id, status)
    VALUES (?, ?, ?, ?, ?, ?, ?)
`

type CreateReportParams struct {
	Url             string
	Title           string
	StartingAt      time.Time
	DurationMinutes int64
	Reporters       string
	ConferenceID    int64
	Status          string
}

func (q *Queries) CreateReport(ctx context.Context, arg CreateReportParams) error {
	_, err := q.exec(ctx, q.createReportStmt, createReport,
		arg.Url,
		arg.Title,
		arg.StartingAt,
		arg.DurationMinutes,
		arg.Reporters,
		arg.ConferenceID,
		arg.Status,
	)
	return err
}

const getAllRatings = `-- name: GetAllRatings :many
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
    report.conference_id = ?
    AND rating.rating_type = 'score'
    AND user.id_data IS NOT NULL
`

type GetAllRatingsRow struct {
	IDData            sql.NullString
	Url               sql.NullString
	ContentScore      sql.NullInt64
	PresentationScore sql.NullInt64
	Notes             sql.NullString
}

func (q *Queries) GetAllRatings(ctx context.Context, conferenceID int64) ([]GetAllRatingsRow, error) {
	rows, err := q.query(ctx, q.getAllRatingsStmt, getAllRatings, conferenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllRatingsRow
	for rows.Next() {
		var i GetAllRatingsRow
		if err := rows.Scan(
			&i.IDData,
			&i.Url,
			&i.ContentScore,
			&i.PresentationScore,
			&i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllReports = `-- name: GetAllReports :many
SELECT
    starting_at,
    duration_minutes,
    title,
    reporters,
    url
FROM report
WHERE conference_id = ?
`

type GetAllReportsRow struct {
	StartingAt      time.Time
	DurationMinutes int64
	Title           string
	Reporters       string
	Url             string
}

func (q *Queries) GetAllReports(ctx context.Context, conferenceID int64) ([]GetAllReportsRow, error) {
	rows, err := q.query(ctx, q.getAllReportsStmt, getAllReports, conferenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllReportsRow
	for rows.Next() {
		var i GetAllReportsRow
		if err := rows.Scan(
			&i.StartingAt,
			&i.DurationMinutes,
			&i.Title,
			&i.Reporters,
			&i.Url,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllUserRatings = `-- name: GetAllUserRatings :many
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
    AND report.conference_id = ?
`

type GetAllUserRatingsParams struct {
	TelegramID   int64
	ConferenceID int64
}

type GetAllUserRatingsRow struct {
	IDData            sql.NullString
	Url               sql.NullString
	ContentScore      sql.NullInt64
	PresentationScore sql.NullInt64
	Notes             sql.NullString
}

func (q *Queries) GetAllUserRatings(ctx context.Context, arg GetAllUserRatingsParams) ([]GetAllUserRatingsRow, error) {
	rows, err := q.query(ctx, q.getAllUserRatingsStmt, getAllUserRatings, arg.TelegramID, arg.ConferenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUserRatingsRow
	for rows.Next() {
		var i GetAllUserRatingsRow
		if err := rows.Scan(
			&i.IDData,
			&i.Url,
			&i.ContentScore,
			&i.PresentationScore,
			&i.Notes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllUserReportsNoScore = `-- name: GetAllUserReportsNoScore :many
SELECT
    starting_at,
    duration_minutes,
    title,
    reporters,
    url
FROM
    report
WHERE
    report.conference_id = ?
    AND id NOT IN (
        SELECT
            report_id
        FROM
            rating
        WHERE
            user_id = ?)
`

type GetAllUserReportsNoScoreParams struct {
	ConferenceID int64
	UserID       int64
}

type GetAllUserReportsNoScoreRow struct {
	StartingAt      time.Time
	DurationMinutes int64
	Title           string
	Reporters       string
	Url             string
}

func (q *Queries) GetAllUserReportsNoScore(ctx context.Context, arg GetAllUserReportsNoScoreParams) ([]GetAllUserReportsNoScoreRow, error) {
	rows, err := q.query(ctx, q.getAllUserReportsNoScoreStmt, getAllUserReportsNoScore, arg.ConferenceID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUserReportsNoScoreRow
	for rows.Next() {
		var i GetAllUserReportsNoScoreRow
		if err := rows.Scan(
			&i.StartingAt,
			&i.DurationMinutes,
			&i.Title,
			&i.Reporters,
			&i.Url,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeFavoriteReport = `-- name: RemoveFavoriteReport :exec
DELETE FROM favorite_reports
WHERE user_id = ?
    AND report_id = ?
`

type RemoveFavoriteReportParams struct {
	UserID   int64
	ReportID int64
}

func (q *Queries) RemoveFavoriteReport(ctx context.Context, arg RemoveFavoriteReportParams) error {
	_, err := q.exec(ctx, q.removeFavoriteReportStmt, removeFavoriteReport, arg.UserID, arg.ReportID)
	return err
}

const updateUserData = `-- name: UpdateUserData :exec
UPDATE
    user
SET
    id_data = ?
WHERE
    telegram_id = ?
`

type UpdateUserDataParams struct {
	IDData     sql.NullString
	TelegramID int64
}

func (q *Queries) UpdateUserData(ctx context.Context, arg UpdateUserDataParams) error {
	_, err := q.exec(ctx, q.updateUserDataStmt, updateUserData, arg.IDData, arg.TelegramID)
	return err
}
