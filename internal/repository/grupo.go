package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type GrupoRepository struct {
	db *sql.DB
}

func NewGrupoRepository(db *sql.DB) *GrupoRepository {
	return &GrupoRepository{db: db}
}

func (r *GrupoRepository) GetAll(offset, limit int, idGrado, idSede, idJornada int) ([]models.Grupo, int, error) {
	query := `SELECT id_grupo, nombre, id_grado, id_sede, id_jornada FROM grupos WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if idGrado > 0 {
		query += ` AND id_grado = $` + string(rune('0'+argNum))
		args = append(args, idGrado)
		argNum++
	}
	if idSede > 0 {
		query += ` AND id_sede = $` + string(rune('0'+argNum))
		args = append(args, idSede)
		argNum++
	}
	if idJornada > 0 {
		query += ` AND id_jornada = $` + string(rune('0'+argNum))
		args = append(args, idJornada)
		argNum++
	}

	query += ` ORDER BY nombre LIMIT $` + string(rune('0'+argNum)) + ` OFFSET $` + string(rune('0'+argNum+1))
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var grupos []models.Grupo
	for rows.Next() {
		var g models.Grupo
		if err := rows.Scan(&g.IDGrupo, &g.Nombre, &g.IDGrado, &g.IDSede, &g.IDJornada); err != nil {
			return nil, 0, err
		}
		grupos = append(grupos, g)
	}

	var total int
	r.db.QueryRow(`SELECT COUNT(*) FROM grupos`).Scan(&total)

	return grupos, total, nil
}

func (r *GrupoRepository) GetByID(id int) (*models.Grupo, error) {
	var g models.Grupo
	err := r.db.QueryRow(
		`SELECT id_grupo, nombre, id_grado, id_sede, id_jornada FROM grupos WHERE id_grupo = $1`,
		id,
	).Scan(&g.IDGrupo, &g.Nombre, &g.IDGrado, &g.IDSede, &g.IDJornada)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GrupoRepository) Create(g *models.Grupo) error {
	return r.db.QueryRow(
		`INSERT INTO grupos (nombre, id_grado, id_sede, id_jornada) VALUES ($1, $2, $3, $4) RETURNING id_grupo`,
		g.Nombre, g.IDGrado, g.IDSede, g.IDJornada,
	).Scan(&g.IDGrupo)
}

func (r *GrupoRepository) Update(g *models.Grupo) error {
	_, err := r.db.Exec(
		`UPDATE grupos SET nombre = $1, id_grado = $2, id_sede = $3, id_jornada = $4 WHERE id_grupo = $5`,
		g.Nombre, g.IDGrado, g.IDSede, g.IDJornada, g.IDGrupo,
	)
	return err
}

func (r *GrupoRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM grupos WHERE id_grupo = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
