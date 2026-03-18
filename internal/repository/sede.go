package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type SedeRepository struct {
	db *sql.DB
}

func NewSedeRepository(db *sql.DB) *SedeRepository {
	return &SedeRepository{db: db}
}

func (r *SedeRepository) GetAll(offset, limit int, idInstitucion int, nombre string) ([]models.Sede, int, error) {
	query := `SELECT id_sede, nombre, direccion, id_institucion, id_municipio 
			  FROM sedes WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM sedes WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if idInstitucion > 0 {
		query += ` AND id_institucion = $` + string(rune('0'+argNum))
		countQuery += ` AND id_institucion = $` + string(rune('0'+argNum))
		args = append(args, idInstitucion)
		argNum++
	}
	if nombre != "" {
		query += ` AND nombre ILIKE $` + string(rune('0'+argNum))
		countQuery += ` AND nombre ILIKE $` + string(rune('0'+argNum))
		args = append(args, "%"+nombre+"%")
		argNum++
	}

	query += ` ORDER BY nombre LIMIT $` + string(rune('0'+argNum)) + ` OFFSET $` + string(rune('0'+argNum+1))
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sedes []models.Sede
	for rows.Next() {
		var s models.Sede
		if err := rows.Scan(&s.IDSede, &s.Nombre, &s.Direccion, &s.IDInstitucion, &s.IDMunicipio); err != nil {
			return nil, 0, err
		}
		sedes = append(sedes, s)
	}

	var total int
	r.db.QueryRow(countQuery).Scan(&total)

	return sedes, total, nil
}

func (r *SedeRepository) GetByID(id int) (*models.Sede, error) {
	var s models.Sede
	err := r.db.QueryRow(
		`SELECT id_sede, nombre, direccion, id_institucion, id_municipio FROM sedes WHERE id_sede = $1`,
		id,
	).Scan(&s.IDSede, &s.Nombre, &s.Direccion, &s.IDInstitucion, &s.IDMunicipio)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SedeRepository) Create(s *models.Sede) error {
	return r.db.QueryRow(
		`INSERT INTO sedes (nombre, direccion, id_institucion, id_municipio) VALUES ($1, $2, $3, $4) 
		 RETURNING id_sede`,
		s.Nombre, s.Direccion, s.IDInstitucion, s.IDMunicipio,
	).Scan(&s.IDSede)
}

func (r *SedeRepository) Update(s *models.Sede) error {
	_, err := r.db.Exec(
		`UPDATE sedes SET nombre = $1, direccion = $2, id_institucion = $3, id_municipio = $4 WHERE id_sede = $5`,
		s.Nombre, s.Direccion, s.IDInstitucion, s.IDMunicipio, s.IDSede,
	)
	return err
}

func (r *SedeRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM sedes WHERE id_sede = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
