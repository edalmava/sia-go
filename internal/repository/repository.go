package repository

import "database/sql"

type Repository struct {
	Institucion  *InstitucionRepository
	Sede         *SedeRepository
	Grado        *GradoRepository
	Grupo        *GrupoRepository
	Jornada      *JornadaRepository
	Asignatura   *AsignaturaRepository
	Docente      *DocenteRepository
	Estudiante   *EstudianteRepository
	Usuario      *UsuarioRepository
	Permiso      *PermisoRepository
	Rol          *RolRepository
	Modulo       *ModuloRepository
	RefreshToken *RefreshTokenRepository
	RevokedToken *RevokedTokenRepository
}

func NewRepository(db *sql.DB) *Repository {
	if db == nil {
		return &Repository{}
	}
	return &Repository{
		Institucion:  NewInstitucionRepository(db),
		Sede:         NewSedeRepository(db),
		Grado:        NewGradoRepository(db),
		Grupo:        NewGrupoRepository(db),
		Jornada:      NewJornadaRepository(db),
		Asignatura:   NewAsignaturaRepository(db),
		Docente:      NewDocenteRepository(db),
		Estudiante:   NewEstudianteRepository(db),
		Usuario:      NewUsuarioRepository(db),
		Permiso:      NewPermisoRepository(db),
		Rol:          NewRolRepository(db),
		Modulo:       NewModuloRepository(db),
		RefreshToken: NewRefreshTokenRepository(db),
		RevokedToken: NewRevokedTokenRepository(db),
	}
}
