package repository

import (
	"time"

	"github.com/edalmava/sia/internal/cache"
	"github.com/edalmava/sia/internal/models"
)

type CachedPermisoRepository struct {
	repo       *PermisoRepository
	byUserID   *cache.Cache[[]string]
	byUserPerm *cache.Cache[bool]
	all        *cache.Cache[[]models.Permiso]
}

func NewCachedPermisoRepository(repo *PermisoRepository) *CachedPermisoRepository {
	return &CachedPermisoRepository{
		repo:       repo,
		byUserID:   cache.New[[]string]("permisos_by_user", 5*time.Minute),
		byUserPerm: cache.New[bool]("permisos_has", 5*time.Minute),
		all:        cache.New[[]models.Permiso]("permisos_all", 30*time.Minute),
	}
}

func (r *CachedPermisoRepository) GetPermissionsByUserID(userID int) ([]string, error) {
	key := r.byUserID.Key(userID)
	if val, ok := r.byUserID.Get(key); ok {
		return val, nil
	}
	val, err := r.repo.GetPermissionsByUserID(userID)
	if err != nil {
		return nil, err
	}
	r.byUserID.Set(key, val)
	return val, nil
}

func (r *CachedPermisoRepository) HasPermission(userID int, permisoCodigo string) (bool, error) {
	key := r.byUserPerm.Key(userID, permisoCodigo)
	if val, ok := r.byUserPerm.Get(key); ok {
		return val, nil
	}
	val, err := r.repo.HasPermission(userID, permisoCodigo)
	if err != nil {
		return false, err
	}
	r.byUserPerm.Set(key, val)
	return val, nil
}

func (r *CachedPermisoRepository) GetAll() ([]models.Permiso, error) {
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

func (r *CachedPermisoRepository) InvalidateUser(userID int) {
	r.byUserID.Delete(r.byUserID.Key(userID))
}

func (r *CachedPermisoRepository) InvalidateAll() {
	r.byUserID.Clear()
	r.byUserPerm.Clear()
	r.all.Clear()
}
