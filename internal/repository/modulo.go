package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type ModuloRepository struct {
	db *sql.DB
}

func NewModuloRepository(db *sql.DB) *ModuloRepository {
	return &ModuloRepository{db: db}
}

func (r *ModuloRepository) GetAll() ([]models.Modulo, error) {
	rows, err := r.db.Query(`
		SELECT m.id_modulo, m.nombre, m.descripcion, m.codigo,
			   COUNT(p.id_permiso) as permisos_count
		FROM modulos m
		LEFT JOIN permisos p ON m.id_modulo = p.id_modulo
		GROUP BY m.id_modulo
		ORDER BY m.codigo
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modulos []models.Modulo
	for rows.Next() {
		var m models.Modulo
		if err := rows.Scan(&m.IDModulo, &m.Nombre, &m.Descripcion, &m.Codigo, &m.PermisosCount); err != nil {
			return nil, err
		}
		modulos = append(modulos, m)
	}

	return modulos, nil
}

func (r *ModuloRepository) GetByID(id int) (*models.Modulo, error) {
	var m models.Modulo
	err := r.db.QueryRow(`
		SELECT id_modulo, nombre, descripcion, codigo
		FROM modulos WHERE id_modulo = $1
	`, id).Scan(&m.IDModulo, &m.Nombre, &m.Descripcion, &m.Codigo)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}
