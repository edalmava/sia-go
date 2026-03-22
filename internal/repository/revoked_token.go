package repository

import (
	"database/sql"
	"time"
)

type RevokedTokenRepository struct {
	db *sql.DB
}

func NewRevokedTokenRepository(db *sql.DB) *RevokedTokenRepository {
	return &RevokedTokenRepository{db: db}
}

func (r *RevokedTokenRepository) Add(jti string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`INSERT INTO revoked_access_tokens (jti, fecha_expiracion) VALUES ($1, $2)
		 ON CONFLICT (jti) DO NOTHING`,
		jti, expiresAt,
	)
	return err
}

func (r *RevokedTokenRepository) IsRevoked(jti string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM revoked_access_tokens WHERE jti = $1)`,
		jti,
	).Scan(&exists)
	return exists, err
}

func (r *RevokedTokenRepository) CleanupExpired() (int64, error) {
	result, err := r.db.Exec(
		`DELETE FROM revoked_access_tokens WHERE fecha_expiracion < $1`,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
