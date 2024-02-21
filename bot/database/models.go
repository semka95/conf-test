// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
	"time"
)

type FavoriteReport struct {
	UserID   int64
	ReportID int64
}

type Rating struct {
	ID                int64
	ReportID          int64
	UserID            int64
	RatingType        string
	ContentScore      sql.NullInt64
	PresentationScore sql.NullInt64
	Notes             sql.NullString
}

type Report struct {
	ID              int64
	Url             string
	Title           string
	StartingAt      time.Time
	DurationMinutes int64
	Reporters       string
	Status          string
}

type User struct {
	TelegramID int64
	IDData     sql.NullString
	Role       string
}
