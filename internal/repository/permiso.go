package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type PermisoRepository struct {
	db *sql.DB
}

func NewPermisoRepository(db *sql.DB) *PermisoRepository {
	return &PermisoRepository{db: db}
}

func (r *PermisoRepository) GetPermissionsByUserID(userID int) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT p.codigo 
		FROM permisos p
		JOIN roles_permisos rp ON p.id_permiso = rp.id_permiso
		JOIN usuarios u ON u.id_rol = rp.id_rol
		WHERE u.id_usuario = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permisos []string
	for rows.Next() {
		var codigo string
		if err := rows.Scan(&codigo); err != nil {
			return nil, err
		}
		permisos = append(permisos, codigo)
	}

	return permisos, nil
}

func (r *PermisoRepository) HasPermission(userID int, permisoCodigo string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM permisos p
			JOIN roles_permisos rp ON p.id_permiso = rp.id_permiso
			JOIN usuarios u ON u.id_rol = rp.id_rol
			WHERE u.id_usuario = $1 AND p.codigo = $2
		)
	`, userID, permisoCodigo).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PermisoRepository) GetAll() ([]models.Permiso, error) {
	rows, err := r.db.Query(`
		SELECT p.id_permiso, p.nombre, p.descripcion, p.codigo, p.id_modulo, m.nombre as modulo_nombre, m.codigo as modulo_codigo
		FROM permisos p
		JOIN modulos m ON p.id_modulo = m.id_modulo
		ORDER BY m.codigo, p.codigo
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permisos []models.Permiso
	for rows.Next() {
		var p models.Permiso
		if err := rows.Scan(&p.IDPermiso, &p.Nombre, &p.Descripcion, &p.Codigo, &p.IDModulo, &p.ModuloNombre, &p.ModuloCodigo); err != nil {
			return nil, err
		}
		permisos = append(permisos, p)
	}

	return permisos, nil
}
