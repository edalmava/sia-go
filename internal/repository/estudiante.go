package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type EstudianteRepository struct {
	db *sql.DB
}

func NewEstudianteRepository(db *sql.DB) *EstudianteRepository {
	return &EstudianteRepository{db: db}
}

func (r *EstudianteRepository) GetAll(offset, limit int) ([]models.Estudiante, int, error) {
	rows, err := r.db.Query(
		`SELECT id_estudiante, documento_identidad, tipo_documento, fecha_nacimiento, telefono, correo_electronico, nombres, apellidos, direccion, id_municipio 
		 FROM estudiantes ORDER BY apellidos, nombres LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var estudiantes []models.Estudiante
	for rows.Next() {
		var e models.Estudiante
		if err := rows.Scan(&e.IDEstudiante, &e.DocumentoIdentidad, &e.TipoDocumento, &e.FechaNacimiento, &e.Telefono, &e.CorreoElectronico, &e.Nombres, &e.Apellidos, &e.Direccion, &e.IDMunicipio); err != nil {
			return nil, 0, err
		}
		estudiantes = append(estudiantes, e)
	}

	var total int
	r.db.QueryRow(`SELECT COUNT(*) FROM estudiantes`).Scan(&total)

	return estudiantes, total, nil
}

func (r *EstudianteRepository) GetByID(id int) (*models.Estudiante, error) {
	var e models.Estudiante
	err := r.db.QueryRow(
		`SELECT id_estudiante, documento_identidad, tipo_documento, fecha_nacimiento, telefono, correo_electronico, nombres, apellidos, direccion, id_municipio 
		 FROM estudiantes WHERE id_estudiante = $1`,
		id,
	).Scan(&e.IDEstudiante, &e.DocumentoIdentidad, &e.TipoDocumento, &e.FechaNacimiento, &e.Telefono, &e.CorreoElectronico, &e.Nombres, &e.Apellidos, &e.Direccion, &e.IDMunicipio)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EstudianteRepository) Create(e *models.Estudiante) error {
	return r.db.QueryRow(
		`INSERT INTO estudiantes (documento_identidad, tipo_documento, fecha_nacimiento, telefono, correo_electronico, nombres, apellidos, direccion, id_municipio) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id_estudiante`,
		e.DocumentoIdentidad, e.TipoDocumento, e.FechaNacimiento, e.Telefono, e.CorreoElectronico, e.Nombres, e.Apellidos, e.Direccion, e.IDMunicipio,
	).Scan(&e.IDEstudiante)
}

func (r *EstudianteRepository) Update(e *models.Estudiante) error {
	_, err := r.db.Exec(
		`UPDATE estudiantes SET documento_identidad = $1, tipo_documento = $2, fecha_nacimiento = $3, telefono = $4, correo_electronico = $5, nombres = $6, apellidos = $7, direccion = $8, id_municipio = $9 
		 WHERE id_estudiante = $10`,
		e.DocumentoIdentidad, e.TipoDocumento, e.FechaNacimiento, e.Telefono, e.CorreoElectronico, e.Nombres, e.Apellidos, e.Direccion, e.IDMunicipio, e.IDEstudiante,
	)
	return err
}

func (r *EstudianteRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM estudiantes WHERE id_estudiante = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
