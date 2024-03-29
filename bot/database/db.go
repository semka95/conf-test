// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.addFavoriteReportStmt, err = db.PrepareContext(ctx, addFavoriteReport); err != nil {
		return nil, fmt.Errorf("error preparing query AddFavoriteReport: %w", err)
	}
	if q.createRatingStmt, err = db.PrepareContext(ctx, createRating); err != nil {
		return nil, fmt.Errorf("error preparing query CreateRating: %w", err)
	}
	if q.createReportStmt, err = db.PrepareContext(ctx, createReport); err != nil {
		return nil, fmt.Errorf("error preparing query CreateReport: %w", err)
	}
	if q.getAllRatingsStmt, err = db.PrepareContext(ctx, getAllRatings); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllRatings: %w", err)
	}
	if q.getAllReportsStmt, err = db.PrepareContext(ctx, getAllReports); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllReports: %w", err)
	}
	if q.getAllUserRatingsStmt, err = db.PrepareContext(ctx, getAllUserRatings); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllUserRatings: %w", err)
	}
	if q.getAllUserReportsNoScoreStmt, err = db.PrepareContext(ctx, getAllUserReportsNoScore); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllUserReportsNoScore: %w", err)
	}
	if q.removeFavoriteReportStmt, err = db.PrepareContext(ctx, removeFavoriteReport); err != nil {
		return nil, fmt.Errorf("error preparing query RemoveFavoriteReport: %w", err)
	}
	if q.updateUserDataStmt, err = db.PrepareContext(ctx, updateUserData); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserData: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addFavoriteReportStmt != nil {
		if cerr := q.addFavoriteReportStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addFavoriteReportStmt: %w", cerr)
		}
	}
	if q.createRatingStmt != nil {
		if cerr := q.createRatingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createRatingStmt: %w", cerr)
		}
	}
	if q.createReportStmt != nil {
		if cerr := q.createReportStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createReportStmt: %w", cerr)
		}
	}
	if q.getAllRatingsStmt != nil {
		if cerr := q.getAllRatingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllRatingsStmt: %w", cerr)
		}
	}
	if q.getAllReportsStmt != nil {
		if cerr := q.getAllReportsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllReportsStmt: %w", cerr)
		}
	}
	if q.getAllUserRatingsStmt != nil {
		if cerr := q.getAllUserRatingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllUserRatingsStmt: %w", cerr)
		}
	}
	if q.getAllUserReportsNoScoreStmt != nil {
		if cerr := q.getAllUserReportsNoScoreStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllUserReportsNoScoreStmt: %w", cerr)
		}
	}
	if q.removeFavoriteReportStmt != nil {
		if cerr := q.removeFavoriteReportStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing removeFavoriteReportStmt: %w", cerr)
		}
	}
	if q.updateUserDataStmt != nil {
		if cerr := q.updateUserDataStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserDataStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                           DBTX
	tx                           *sql.Tx
	addFavoriteReportStmt        *sql.Stmt
	createRatingStmt             *sql.Stmt
	createReportStmt             *sql.Stmt
	getAllRatingsStmt            *sql.Stmt
	getAllReportsStmt            *sql.Stmt
	getAllUserRatingsStmt        *sql.Stmt
	getAllUserReportsNoScoreStmt *sql.Stmt
	removeFavoriteReportStmt     *sql.Stmt
	updateUserDataStmt           *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                           tx,
		tx:                           tx,
		addFavoriteReportStmt:        q.addFavoriteReportStmt,
		createRatingStmt:             q.createRatingStmt,
		createReportStmt:             q.createReportStmt,
		getAllRatingsStmt:            q.getAllRatingsStmt,
		getAllReportsStmt:            q.getAllReportsStmt,
		getAllUserRatingsStmt:        q.getAllUserRatingsStmt,
		getAllUserReportsNoScoreStmt: q.getAllUserReportsNoScoreStmt,
		removeFavoriteReportStmt:     q.removeFavoriteReportStmt,
		updateUserDataStmt:           q.updateUserDataStmt,
	}
}
