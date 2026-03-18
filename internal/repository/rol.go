package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type RolRepository struct {
	db *sql.DB
}

func NewRolRepository(db *sql.DB) *RolRepository {
	return &RolRepository{db: db}
}

func (r *RolRepository) GetAll() ([]models.Rol, error) {
	rows, err := r.db.Query(`
		SELECT r.id_rol, r.nombre, r.descripcion, r.es_rol_sistema,
			   COUNT(rp.id_permiso) as permisos_count
		FROM roles r
		LEFT JOIN roles_permisos rp ON r.id_rol = rp.id_rol
		GROUP BY r.id_rol
		ORDER BY r.id_rol
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Rol
	for rows.Next() {
		var rol models.Rol
		if err := rows.Scan(&rol.IDRol, &rol.Nombre, &rol.Descripcion, &rol.EsRolSistema, &rol.PermisosCount); err != nil {
			return nil, err
		}
		roles = append(roles, rol)
	}

	return roles, nil
}

func (r *RolRepository) GetByID(id int) (*models.Rol, error) {
	var rol models.Rol
	err := r.db.QueryRow(`
		SELECT id_rol, nombre, descripcion, es_rol_sistema
		FROM roles WHERE id_rol = $1
	`, id).Scan(&rol.IDRol, &rol.Nombre, &rol.Descripcion, &rol.EsRolSistema)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rol, nil
}

func (r *RolRepository) GetByIDWithPermisos(id int) (*models.Rol, error) {
	rol, err := r.GetByID(id)
	if err != nil || rol == nil {
		return rol, err
	}

	rows, err := r.db.Query(`
		SELECT p.id_permiso, p.nombre, p.descripcion, p.codigo, p.id_modulo
		FROM permisos p
		JOIN roles_permisos rp ON p.id_permiso = rp.id_permiso
		WHERE rp.id_rol = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permisos []models.Permiso
	for rows.Next() {
		var p models.Permiso
		if err := rows.Scan(&p.IDPermiso, &p.Nombre, &p.Descripcion, &p.Codigo, &p.IDModulo); err != nil {
			return nil, err
		}
		permisos = append(permisos, p)
	}

	rol.Permisos = permisos
	return rol, nil
}

func (r *RolRepository) Create(rol *models.Rol) error {
	return r.db.QueryRow(`
		INSERT INTO roles (nombre, descripcion, es_rol_sistema)
		VALUES ($1, $2, $3)
		RETURNING id_rol
	`, rol.Nombre, rol.Descripcion, rol.EsRolSistema).Scan(&rol.IDRol)
}

func (r *RolRepository) Update(rol *models.Rol) error {
	_, err := r.db.Exec(`
		UPDATE roles SET nombre = $1, descripcion = $2, es_rol_sistema = $3
		WHERE id_rol = $4
	`, rol.Nombre, rol.Descripcion, rol.EsRolSistema, rol.IDRol)
	return err
}

func (r *RolRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM roles WHERE id_rol = $1 AND es_rol_sistema = FALSE`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *RolRepository) SetPermisos(rolID int, permisos []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM roles_permisos WHERE id_rol = $1`, rolID)
	if err != nil {
		return err
	}

	for _, permisoID := range permisos {
		_, err = tx.Exec(`
			INSERT INTO roles_permisos (id_rol, id_permiso) VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, rolID, permisoID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
