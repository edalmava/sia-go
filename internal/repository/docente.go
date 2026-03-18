package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type DocenteRepository struct {
	db *sql.DB
}

func NewDocenteRepository(db *sql.DB) *DocenteRepository {
	return &DocenteRepository{db: db}
}

func (r *DocenteRepository) GetAll(offset, limit int) ([]models.Docente, int, error) {
	rows, err := r.db.Query(
		`SELECT id_docente, nombres, apellidos, documento_identidad, tipo_documento, profesion, telefono, correo_electronico 
		 FROM docentes ORDER BY apellidos, nombres LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var docentes []models.Docente
	for rows.Next() {
		var d models.Docente
		if err := rows.Scan(&d.IDDocente, &d.Nombres, &d.Apellidos, &d.DocumentoIdentidad, &d.TipoDocumento, &d.Profesion, &d.Telefono, &d.CorreoElectronico); err != nil {
			return nil, 0, err
		}
		docentes = append(docentes, d)
	}

	var total int
	r.db.QueryRow(`SELECT COUNT(*) FROM docentes`).Scan(&total)

	return docentes, total, nil
}

func (r *DocenteRepository) GetByID(id int) (*models.Docente, error) {
	var d models.Docente
	err := r.db.QueryRow(
		`SELECT id_docente, nombres, apellidos, documento_identidad, tipo_documento, profesion, telefono, correo_electronico 
		 FROM docentes WHERE id_docente = $1`,
		id,
	).Scan(&d.IDDocente, &d.Nombres, &d.Apellidos, &d.DocumentoIdentidad, &d.TipoDocumento, &d.Profesion, &d.Telefono, &d.CorreoElectronico)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DocenteRepository) Create(d *models.Docente) error {
	return r.db.QueryRow(
		`INSERT INTO docentes (nombres, apellidos, documento_identidad, tipo_documento, profesion, telefono, correo_electronico) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_docente`,
		d.Nombres, d.Apellidos, d.DocumentoIdentidad, d.TipoDocumento, d.Profesion, d.Telefono, d.CorreoElectronico,
	).Scan(&d.IDDocente)
}

func (r *DocenteRepository) Update(d *models.Docente) error {
	_, err := r.db.Exec(
		`UPDATE docentes SET nombres = $1, apellidos = $2, documento_identidad = $3, tipo_documento = $4, profesion = $5, telefono = $6, correo_electronico = $7 
		 WHERE id_docente = $8`,
		d.Nombres, d.Apellidos, d.DocumentoIdentidad, d.TipoDocumento, d.Profesion, d.Telefono, d.CorreoElectronico, d.IDDocente,
	)
	return err
}

func (r *DocenteRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM docentes WHERE id_docente = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
