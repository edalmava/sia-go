package repository

import (
	"database/sql"

	"github.com/edalmava/sia/internal/models"
)

type GradoRepository struct {
	db *sql.DB
}

func NewGradoRepository(db *sql.DB) *GradoRepository {
	return &GradoRepository{db: db}
}

func (r *GradoRepository) GetAll(offset, limit int, nombre string) ([]models.Grado, int, error) {
	query := `SELECT id_grado, nombre FROM grados WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM grados WHERE 1=1`
	args := []interface{}{}

	if nombre != "" {
		query += ` AND nombre ILIKE $1`
		countQuery += ` AND nombre ILIKE $1`
		args = append(args, "%"+nombre+"%")
	}

	query += ` ORDER BY id_grado LIMIT $` + string(rune('0'+len(args)+1)) + ` OFFSET $` + string(rune('0'+len(args)+2))
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var grados []models.Grado
	for rows.Next() {
		var g models.Grado
		if err := rows.Scan(&g.IDGrado, &g.Nombre); err != nil {
			return nil, 0, err
		}
		grados = append(grados, g)
	}

	var total int
	r.db.QueryRow(countQuery).Scan(&total)

	return grados, total, nil
}

func (r *GradoRepository) GetByID(id int) (*models.Grado, error) {
	var g models.Grado
	err := r.db.QueryRow(`SELECT id_grado, nombre FROM grados WHERE id_grado = $1`, id).Scan(&g.IDGrado, &g.Nombre)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GradoRepository) Create(g *models.Grado) error {
	return r.db.QueryRow(`INSERT INTO grados (nombre) VALUES ($1) RETURNING id_grado`, g.Nombre).Scan(&g.IDGrado)
}

func (r *GradoRepository) Update(g *models.Grado) error {
	_, err := r.db.Exec(`UPDATE grados SET nombre = $1 WHERE id_grado = $2`, g.Nombre, g.IDGrado)
	return err
}

func (r *GradoRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM grados WHERE id_grado = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
