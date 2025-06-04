package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// TrackingInfo represents email tracking data.
type TrackingInfo struct {
	ID        string
	Email     string
	Campaign  string
	UTMSource string
	UTMMedium string
	Clicked   bool
	CreatedAt time.Time
	ClickedAt sql.NullTime
}

// Repository handles persistence of tracking info in Postgres.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new Repository with the given connection string.
func NewRepository(conn string) (*Repository, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Close closes the database connection.
func (r *Repository) Close() error { return r.db.Close() }

// Create inserts a new tracking record and returns its ID.
func (r *Repository) Create(ctx context.Context, info *TrackingInfo) (string, error) {
	query := `INSERT INTO trackings
    (email, campaign, utm_source, utm_medium, clicked, created_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		info.Email, info.Campaign, info.UTMSource, info.UTMMedium,
		info.Clicked, info.CreatedAt).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// MarkClicked sets the clicked flag and time for the given ID.
func (r *Repository) MarkClicked(ctx context.Context, id string) error {
	query := `UPDATE trackings SET clicked = true, clicked_at = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, time.Now())
	return err
}

// Get retrieves tracking info by ID.
func (r *Repository) Get(ctx context.Context, id string) (*TrackingInfo, error) {
	query := `SELECT id, email, campaign, utm_source, utm_medium, clicked,
        created_at, clicked_at FROM trackings WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var info TrackingInfo
	if err := row.Scan(&info.ID, &info.Email, &info.Campaign, &info.UTMSource,
		&info.UTMMedium, &info.Clicked, &info.CreatedAt, &info.ClickedAt); err != nil {
		return nil, err
	}
	return &info, nil
}
