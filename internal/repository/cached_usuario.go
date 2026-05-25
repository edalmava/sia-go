package repository

import (
	"time"

	"github.com/edalmava/sia/internal/cache"
	"github.com/edalmava/sia/internal/models"
)

type CachedUsuarioRepository struct {
	repo *UsuarioRepository
	all  *cache.Cache[struct {
		usuarios []models.UsuarioResponse
		total    int
	}]
	byID       *cache.Cache[*models.UsuarioResponse]
	byUsername *cache.Cache[*models.Usuario]
}

func NewCachedUsuarioRepository(repo *UsuarioRepository) *CachedUsuarioRepository {
	return &CachedUsuarioRepository{
		repo: repo,
		all: cache.New[struct {
			usuarios []models.UsuarioResponse
			total    int
		}]("usuarios_all", 2*time.Minute),
		byID:       cache.New[*models.UsuarioResponse]("usuario_by_id", 2*time.Minute),
		byUsername: cache.New[*models.Usuario]("usuario_by_username", 2*time.Minute),
	}
}

func (r *CachedUsuarioRepository) GetByUsername(username string) (*models.Usuario, error) {
	key := r.byUsername.Key(username)
	if val, ok := r.byUsername.Get(key); ok {
		return val, nil
	}
	val, err := r.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	r.byUsername.Set(key, val)
	return val, nil
}

func (r *CachedUsuarioRepository) GetByRol(rolID int) ([]models.UsuarioResponse, error) {
	return r.repo.GetByRol(rolID)
}

func (r *CachedUsuarioRepository) GetByID(id int) (*models.UsuarioResponse, error) {
	key := r.byID.Key(id)
	if val, ok := r.byID.Get(key); ok {
		return val, nil
	}
	val, err := r.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	r.byID.Set(key, val)
	return val, nil
}

func (r *CachedUsuarioRepository) GetAll(offset, limit int) ([]models.UsuarioResponse, int, error) {
	key := r.all.Key(offset, limit)
	if val, ok := r.all.Get(key); ok {
		return val.usuarios, val.total, nil
	}
	usuarios, total, err := r.repo.GetAll(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	r.all.Set(key, struct {
		usuarios []models.UsuarioResponse
		total    int
	}{usuarios: usuarios, total: total})
	return usuarios, total, nil
}

func (r *CachedUsuarioRepository) Create(u *models.Usuario) error {
	if err := r.repo.Create(u); err != nil {
		return err
	}
	r.InvalidateAll()
	return nil
}

func (r *CachedUsuarioRepository) Update(u *models.Usuario) error {
	if err := r.repo.Update(u); err != nil {
		return err
	}
	r.InvalidateAll()
	return nil
}

func (r *CachedUsuarioRepository) Delete(id int) error {
	if err := r.repo.Delete(id); err != nil {
		return err
	}
	r.InvalidateAll()
	return nil
}

func (r *CachedUsuarioRepository) UpdatePassword(id int, newPassword string) error {
	if err := r.repo.UpdatePassword(id, newPassword); err != nil {
		return err
	}
	r.InvalidateUser(id)
	return nil
}

func (r *CachedUsuarioRepository) InvalidateUser(id int) {
	r.byID.Delete(r.byID.Key(id))
}

func (r *CachedUsuarioRepository) InvalidateAll() {
	r.all.Clear()
	r.byID.Clear()
	r.byUsername.Clear()
}
