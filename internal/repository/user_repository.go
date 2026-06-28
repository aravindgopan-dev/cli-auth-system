package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID             int
	Username       string
	PasswordHash   string
	TwoFASecret    string
	TwoFAEnabled   bool
	FailedAttempts int
	LockedUntil    sql.NullTime // Maps natively to database null states
	CreatedAt      time.Time
	LastLogin      sql.NullTime // Maps natively to database null states
}

type Session struct {
	Token     string
	Username  string
	ExpiresAt time.Time
}

// Concrete repository struct. No interface declared here!
type Repo struct {
	Conn *pgx.Conn
}

func NewRepo(conn *pgx.Conn) *Repo {
	return &Repo{Conn: conn}
}

func (r *Repo) CreateUser(ctx context.Context, username, passwordHash string) error {
	_, err := r.Conn.Exec(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, passwordHash)
	return err
}

func (r *Repo) GetUser(ctx context.Context, username string) (*User, error) {
	row := r.Conn.QueryRow(ctx, "SELECT id, username, password_hash, two_fa_secret, two_fa_enabled, failed_attempts, locked_until, created_at, last_login FROM users WHERE username = $1", username)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.TwoFASecret, &u.TwoFAEnabled, &u.FailedAttempts, &u.LockedUntil, &u.CreatedAt, &u.LastLogin)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) UpdateUserAuth(ctx context.Context, u *User) error {
	_, err := r.Conn.Exec(ctx, `
		UPDATE users 
		SET password_hash=$1, two_fa_secret=$2, two_fa_enabled=$3, failed_attempts=$4, locked_until=$5, last_login=$6 
		WHERE id=$7`,
		u.PasswordHash, u.TwoFASecret, u.TwoFAEnabled, u.FailedAttempts, u.LockedUntil, u.LastLogin, u.ID)
	return err
}

func (r *Repo) CreateSession(ctx context.Context, token, username string, expiresAt time.Time) error {
	_, err := r.Conn.Exec(ctx, "INSERT INTO sessions (token, username, expires_at) VALUES ($1, $2, $3)", token, username, expiresAt)
	return err
}

func (r *Repo) GetSession(ctx context.Context, token string) (*Session, error) {
	row := r.Conn.QueryRow(ctx, "SELECT token, username, expires_at FROM sessions WHERE token = $1", token)
	var s Session
	err := row.Scan(&s.Token, &s.Username, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repo) DeleteSession(ctx context.Context, token string) error {
	_, err := r.Conn.Exec(ctx, "DELETE FROM sessions WHERE token = $1", token)
	return err
}