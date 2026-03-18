package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type JornadaRepository struct {
	db *sql.DB
}

func NewJornadaRepository(db *sql.DB) *JornadaRepository {
	return &JornadaRepository{db: db}
}

func (r *JornadaRepository) GetAll(offset, limit int, nombre string) ([]models.Jornada, int, error) {
	query := `SELECT id_jornada, nombre FROM jornadas WHERE 1=1`
	args := []interface{}{limit, offset}

	if nombre != "" {
		query = `SELECT id_jornada, nombre FROM jornadas WHERE nombre ILIKE $1 LIMIT $2 OFFSET $3`
		args = []interface{}{"%" + nombre + "%", limit, offset}
	} else {
		query += ` ORDER BY nombre LIMIT $1 OFFSET $2`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jornadas []models.Jornada
	for rows.Next() {
		var j models.Jornada
		if err := rows.Scan(&j.IDJornada, &j.Nombre); err != nil {
			return nil, 0, err
		}
		jornadas = append(jornadas, j)
	}

	var total int
	r.db.QueryRow(`SELECT COUNT(*) FROM jornadas`).Scan(&total)

	return jornadas, total, nil
}

func (r *JornadaRepository) GetByID(id int) (*models.Jornada, error) {
	var j models.Jornada
	err := r.db.QueryRow(`SELECT id_jornada, nombre FROM jornadas WHERE id_jornada = $1`, id).Scan(&j.IDJornada, &j.Nombre)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *JornadaRepository) Create(j *models.Jornada) error {
	return r.db.QueryRow(`INSERT INTO jornadas (nombre) VALUES ($1) RETURNING id_jornada`, j.Nombre).Scan(&j.IDJornada)
}

func (r *JornadaRepository) Update(j *models.Jornada) error {
	_, err := r.db.Exec(`UPDATE jornadas SET nombre = $1 WHERE id_jornada = $2`, j.Nombre, j.IDJornada)
	return err
}

func (r *JornadaRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM jornadas WHERE id_jornada = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
