package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type AsignaturaRepository struct {
	db *sql.DB
}

func NewAsignaturaRepository(db *sql.DB) *AsignaturaRepository {
	return &AsignaturaRepository{db: db}
}

func (r *AsignaturaRepository) GetAll(offset, limit int) ([]models.Asignatura, int, error) {
	rows, err := r.db.Query(`SELECT id_asignatura, nombre, intensidad_horaria FROM asignaturas ORDER BY nombre LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var asignaturas []models.Asignatura
	for rows.Next() {
		var a models.Asignatura
		if err := rows.Scan(&a.IDAsignatura, &a.Nombre, &a.IntensidadHoraria); err != nil {
			return nil, 0, err
		}
		asignaturas = append(asignaturas, a)
	}

	var total int
	r.db.QueryRow(`SELECT COUNT(*) FROM asignaturas`).Scan(&total)

	return asignaturas, total, nil
}

func (r *AsignaturaRepository) GetByID(id int) (*models.Asignatura, error) {
	var a models.Asignatura
	err := r.db.QueryRow(`SELECT id_asignatura, nombre, intensidad_horaria FROM asignaturas WHERE id_asignatura = $1`, id).Scan(&a.IDAsignatura, &a.Nombre, &a.IntensidadHoraria)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AsignaturaRepository) Create(a *models.Asignatura) error {
	return r.db.QueryRow(`INSERT INTO asignaturas (nombre, intensidad_horaria) VALUES ($1, $2) RETURNING id_asignatura`, a.Nombre, a.IntensidadHoraria).Scan(&a.IDAsignatura)
}

func (r *AsignaturaRepository) Update(a *models.Asignatura) error {
	_, err := r.db.Exec(`UPDATE asignaturas SET nombre = $1, intensidad_horaria = $2 WHERE id_asignatura = $3`, a.Nombre, a.IntensidadHoraria, a.IDAsignatura)
	return err
}

func (r *AsignaturaRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM asignaturas WHERE id_asignatura = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
