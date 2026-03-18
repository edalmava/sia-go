package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type UsuarioRepository struct {
	db *sql.DB
}

func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

func (r *UsuarioRepository) GetByUsername(username string) (*models.Usuario, error) {
	var u models.Usuario
	err := r.db.QueryRow(
		`SELECT id_usuario, nombre_usuario, clave, id_docente, id_estudiante, activo, id_rol 
		 FROM usuarios WHERE nombre_usuario = $1 AND activo = true`,
		username,
	).Scan(&u.IDUsuario, &u.NombreUsuario, &u.Clave, &u.IDDocente, &u.IDEstudiante, &u.Activo, &u.IDRol)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsuarioRepository) GetAll(offset, limit int) ([]models.Usuario, int, error) {
	rows, err := r.db.Query(
		`SELECT id_usuario, nombre_usuario, clave, id_docente, id_estudiante, activo, id_rol 
		 FROM usuarios ORDER BY nombre_usuario LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var u models.Usuario
		if err := rows.Scan(&u.IDUsuario, &u.NombreUsuario, &u.Clave, &u.IDDocente, &u.IDEstudiante, &u.Activo, &u.IDRol); err != nil {
			return nil, 0, err
		}
		usuarios = append(usuarios, u)
	}

	var total int
	r.db.QueryRow(`SELECT COUNT(*) FROM usuarios`).Scan(&total)

	return usuarios, total, nil
}

func (r *UsuarioRepository) Create(u *models.Usuario) error {
	return r.db.QueryRow(
		`INSERT INTO usuarios (nombre_usuario, clave, id_docente, id_estudiante, activo, id_rol) 
		 VALUES ($1, $2, $3, $4, $5, $6) 
		 RETURNING id_usuario`,
		u.NombreUsuario, u.Clave, u.IDDocente, u.IDEstudiante, u.Activo, u.IDRol,
	).Scan(&u.IDUsuario)
}

func (r *UsuarioRepository) GetByID(id int) (*models.Usuario, error) {
	var u models.Usuario
	err := r.db.QueryRow(
		`SELECT id_usuario, nombre_usuario, clave, id_docente, id_estudiante, activo, id_rol 
		 FROM usuarios WHERE id_usuario = $1`,
		id,
	).Scan(&u.IDUsuario, &u.NombreUsuario, &u.Clave, &u.IDDocente, &u.IDEstudiante, &u.Activo, &u.IDRol)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsuarioRepository) UpdatePassword(id int, newPassword string) error {
	_, err := r.db.Exec(
		`UPDATE usuarios SET clave = $1 WHERE id_usuario = $2`,
		newPassword, id,
	)
	return err
}

func (r *UsuarioRepository) Update(u *models.Usuario) error {
	_, err := r.db.Exec(
		`UPDATE usuarios SET nombre_usuario = COALESCE(NULLIF($1, ''), nombre_usuario), 
		 id_docente = $2, id_estudiante = $3, activo = $4, id_rol = COALESCE(NULLIF($5, 0), id_rol) 
		 WHERE id_usuario = $6`,
		u.NombreUsuario, u.IDDocente, u.IDEstudiante, u.Activo, u.IDRol, u.IDUsuario,
	)
	return err
}

func (r *UsuarioRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM usuarios WHERE id_usuario = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
