package repository

import (
	"time"

	"github.com/edalmava/sia/internal/cache"
	"github.com/edalmava/sia/internal/models"
)

type CachedModuloRepository struct {
	repo *ModuloRepository
	all  *cache.Cache[[]models.Modulo]
	byID *cache.Cache[*models.Modulo]
}

func NewCachedModuloRepository(repo *ModuloRepository) *CachedModuloRepository {
	return &CachedModuloRepository{
		repo: repo,
		all:  cache.New[[]models.Modulo]("modulos_all", 30*time.Minute),
		byID: cache.New[*models.Modulo]("modulo_by_id", 30*time.Minute),
	}
}

func (r *CachedModuloRepository) GetAll() ([]models.Modulo, error) {
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

func (r *CachedModuloRepository) GetByID(id int) (*models.Modulo, error) {
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

func (r *CachedModuloRepository) InvalidateAll() {
	r.all.Clear()
	r.byID.Clear()
}
