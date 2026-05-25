package repository

import "github.com/edalmava/sia/internal/models"

type UsuarioReader interface {
	GetByUsername(username string) (*models.Usuario, error)
	GetByID(id int) (*models.UsuarioResponse, error)
	GetAll(offset, limit int) ([]models.UsuarioResponse, int, error)
	GetByRol(rolID int) ([]models.UsuarioResponse, error)
}

type UsuarioWriter interface {
	Create(u *models.Usuario) error
	Update(u *models.Usuario) error
	Delete(id int) error
	UpdatePassword(id int, newPassword string) error
}

type PermisoReader interface {
	GetPermissionsByUserID(userID int) ([]string, error)
	HasPermission(userID int, permisoCodigo string) (bool, error)
	GetAll() ([]models.Permiso, error)
}

type RolReader interface {
	GetAll() ([]models.Rol, error)
	GetByID(id int) (*models.Rol, error)
	GetByIDWithPermisos(id int) (*models.Rol, error)
}

type RolWriter interface {
	Create(rol *models.Rol) error
	Update(rol *models.Rol) error
	Delete(id int) error
	SetPermisos(rolID int, permisos []int) error
}

type ModuloReader interface {
	GetAll() ([]models.Modulo, error)
	GetByID(id int) (*models.Modulo, error)
}
