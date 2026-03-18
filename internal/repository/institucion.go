package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type InstitucionRepository struct {
	db *sql.DB
}

func NewInstitucionRepository(db *sql.DB) *InstitucionRepository {
	return &InstitucionRepository{db: db}
}

func (r *InstitucionRepository) GetAll(offset, limit int, nombre string) ([]models.Institucion, int, error) {
	var query string
	var countQuery string
	var args []interface{}

	if nombre != "" {
		query = `SELECT id_institucion, nombre, codigo_dane FROM instituciones 
				 WHERE nombre ILIKE $1 ORDER BY nombre LIMIT $2 OFFSET $3`
		countQuery = `SELECT COUNT(*) FROM instituciones WHERE nombre ILIKE $1`
		args = []interface{}{"%" + nombre + "%", limit, offset}
	} else {
		query = `SELECT id_institucion, nombre, codigo_dane FROM instituciones 
				 ORDER BY nombre LIMIT $1 OFFSET $2`
		countQuery = `SELECT COUNT(*) FROM instituciones`
		args = []interface{}{limit, offset}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var instituciones []models.Institucion
	for rows.Next() {
		var inst models.Institucion
		if err := rows.Scan(&inst.IDInstitucion, &inst.Nombre, &inst.CodigoDane); err != nil {
			return nil, 0, err
		}
		instituciones = append(instituciones, inst)
	}

	var total int
	if nombre != "" {
		r.db.QueryRow(countQuery, "%"+nombre+"%").Scan(&total)
	} else {
		r.db.QueryRow(countQuery).Scan(&total)
	}

	return instituciones, total, nil
}

func (r *InstitucionRepository) GetByID(id int) (*models.Institucion, error) {
	var inst models.Institucion
	err := r.db.QueryRow(
		`SELECT id_institucion, nombre, codigo_dane FROM instituciones WHERE id_institucion = $1`,
		id,
	).Scan(&inst.IDInstitucion, &inst.Nombre, &inst.CodigoDane)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &inst, nil
}

func (r *InstitucionRepository) Create(inst *models.Institucion) error {
	return r.db.QueryRow(
		`INSERT INTO instituciones (nombre, codigo_dane) VALUES ($1, $2) 
		 RETURNING id_institucion`,
		inst.Nombre, inst.CodigoDane,
	).Scan(&inst.IDInstitucion)
}

func (r *InstitucionRepository) Update(inst *models.Institucion) error {
	_, err := r.db.Exec(
		`UPDATE instituciones SET nombre = $1, codigo_dane = $2 WHERE id_institucion = $3`,
		inst.Nombre, inst.CodigoDane, inst.IDInstitucion,
	)
	return err
}

func (r *InstitucionRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM instituciones WHERE id_institucion = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
