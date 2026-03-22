package repository

import (
	"database/sql"
	"time"

	"github.com/edalmava/sia/internal/models"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(tokenHash, jti string, idUsuario int, fechaExpiracion time.Time, dispositivo string) error {
	_, err := r.db.Exec(
		`INSERT INTO refresh_tokens (token_hash, jti, id_usuario, fecha_expiracion, dispositivo, activo) 
		 VALUES ($1, $2, $3, $4, $5, true)`,
		tokenHash, jti, idUsuario, fechaExpiracion, sql.NullString{String: dispositivo, Valid: dispositivo != ""},
	)
	return err
}

func (r *RefreshTokenRepository) GetByTokenHash(tokenHash string) (*models.RefreshTokenDB, error) {
	var rt models.RefreshTokenDB
	err := r.db.QueryRow(
		`SELECT id, token_hash, jti, id_usuario, fecha_expiracion, fecha_creacion, dispositivo, activo 
		 FROM refresh_tokens WHERE token_hash = $1 AND activo = true`,
		tokenHash,
	).Scan(&rt.ID, &rt.TokenHash, &rt.JTI, &rt.IDUsuario, &rt.FechaExpiracion, &rt.FechaCreacion, &rt.Dispositivo, &rt.Activo)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) Revoke(tokenHash string) error {
	_, err := r.db.Exec(
		`UPDATE refresh_tokens SET activo = false WHERE token_hash = $1`,
		tokenHash,
	)
	return err
}

func (r *RefreshTokenRepository) RevokeAllForUser(idUsuario int) error {
	_, err := r.db.Exec(
		`UPDATE refresh_tokens SET activo = false WHERE id_usuario = $1`,
		idUsuario,
	)
	return err
}

func (r *RefreshTokenRepository) RevokeAll() error {
	_, err := r.db.Exec(`UPDATE refresh_tokens SET activo = false`)
	return err
}

func (r *RefreshTokenRepository) DeleteExpired() (int64, error) {
	result, err := r.db.Exec(
		`DELETE FROM refresh_tokens WHERE fecha_expiracion < $1 OR activo = false`,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
