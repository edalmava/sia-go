package repository

import (
	"time"

	"github.com/edalmava/sia/internal/cache"
	"github.com/edalmava/sia/internal/models"
)

type CachedRolRepository struct {
	repo           *RolRepository
	all            *cache.Cache[[]models.Rol]
	byID           *cache.Cache[*models.Rol]
	byIDWithPermis *cache.Cache[*models.Rol]
}

func NewCachedRolRepository(repo *RolRepository) *CachedRolRepository {
	return &CachedRolRepository{
		repo:           repo,
		all:            cache.New[[]models.Rol]("roles_all", 10*time.Minute),
		byID:           cache.New[*models.Rol]("rol_by_id", 10*time.Minute),
		byIDWithPermis: cache.New[*models.Rol]("rol_with_permisos", 10*time.Minute),
	}
}

func (r *CachedRolRepository) GetAll() ([]models.Rol, error) {
	key := r.all.Key()
	if val, ok := r.all.Get(key); ok {
		return val, nil
	}
	val, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}
	r.all.Set(key, val)
	return val, nil
}

func (r *CachedRolRepository) GetByID(id int) (*models.Rol, error) {
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

func (r *CachedRolRepository) GetByIDWithPermisos(id int) (*models.Rol, error) {
	key := r.byIDWithPermis.Key(id)
	if val, ok := r.byIDWithPermis.Get(key); ok {
		return val, nil
	}
	val, err := r.repo.GetByIDWithPermisos(id)
	if err != nil {
		return nil, err
	}
	r.byIDWithPermis.Set(key, val)
	return val, nil
}

func (r *CachedRolRepository) Create(rol *models.Rol) error {
	if err := r.repo.Create(rol); err != nil {
		return err
	}
	r.InvalidateAll()
	return nil
}

func (r *CachedRolRepository) Update(rol *models.Rol) error {
	if err := r.repo.Update(rol); err != nil {
		return err
	}
	r.InvalidateAll()
	return nil
}

func (r *CachedRolRepository) Delete(id int) error {
	if err := r.repo.Delete(id); err != nil {
		return err
	}
	r.InvalidateAll()
	return nil
}

func (r *CachedRolRepository) SetPermisos(rolID int, permisos []int) error {
	if err := r.repo.SetPermisos(rolID, permisos); err != nil {
		return err
	}
	r.InvalidateAll()
	return nil
}

func (r *CachedRolRepository) InvalidateAll() {
	r.all.Clear()
	r.byID.Clear()
	r.byIDWithPermis.Clear()
}
