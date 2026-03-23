package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/edalmava/sia/internal/middleware"
	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/edalmava/sia/internal/utils"
	"github.com/labstack/echo/v4"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

func parseIDParam(c echo.Context, paramName string) (int, error) {
	id, err := strconv.Atoi(c.Param(paramName))
	if err != nil {
		return 0, c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}
	return id, nil
}

func parsePagination(c echo.Context) (int, int) {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return offset, limit
}

type EstudianteHandler struct {
	repo *repository.EstudianteRepository
}

func NewEstudianteHandler(repo *repository.EstudianteRepository) *EstudianteHandler {
	return &EstudianteHandler{repo: repo}
}

func (h *EstudianteHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	offset, limit := parsePagination(c)
	estudiantes, total, err := h.repo.GetAll(offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener estudiantes",
		})
	}
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       estudiantes,
		Pagination: models.Pagination{Total: total, Offset: offset, Limit: limit},
	})
}

func (h *EstudianteHandler) GetByID(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	e, err := h.repo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener estudiante",
		})
	}
	if e == nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Estudiante no encontrado",
		})
	}
	return c.JSON(http.StatusOK, e)
}

func (h *EstudianteHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	var e models.Estudiante
	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	if err := h.repo.Create(&e); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear estudiante",
		})
	}
	return c.JSON(http.StatusCreated, e)
}

func (h *EstudianteHandler) Update(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var e models.Estudiante
	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	e.IDEstudiante = id
	if err := h.repo.Update(&e); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar estudiante",
		})
	}
	return c.JSON(http.StatusOK, e)
}

func (h *EstudianteHandler) Delete(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.repo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al eliminar estudiante",
		})
	}
	return c.NoContent(http.StatusNoContent)
}

type DocenteHandler struct {
	repo *repository.DocenteRepository
}

func NewDocenteHandler(repo *repository.DocenteRepository) *DocenteHandler {
	return &DocenteHandler{repo: repo}
}

func (h *DocenteHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	offset, limit := parsePagination(c)
	docentes, total, err := h.repo.GetAll(offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener docentes",
		})
	}
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       docentes,
		Pagination: models.Pagination{Total: total, Offset: offset, Limit: limit},
	})
}

func (h *DocenteHandler) GetByID(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	d, err := h.repo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener docente",
		})
	}
	if d == nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Docente no encontrado",
		})
	}
	return c.JSON(http.StatusOK, d)
}

func (h *DocenteHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	var d models.Docente
	if err := c.Bind(&d); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	if err := h.repo.Create(&d); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear docente",
		})
	}
	return c.JSON(http.StatusCreated, d)
}

func (h *DocenteHandler) Update(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var d models.Docente
	if err := c.Bind(&d); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	d.IDDocente = id
	if err := h.repo.Update(&d); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar docente",
		})
	}
	return c.JSON(http.StatusOK, d)
}

func (h *DocenteHandler) Delete(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.repo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al eliminar docente",
		})
	}
	return c.NoContent(http.StatusNoContent)
}

type MatriculaHandler struct{}

func NewMatriculaHandler() *MatriculaHandler {
	return &MatriculaHandler{}
}

func (h *MatriculaHandler) GetAll(c echo.Context) error {
	offset, limit := parsePagination(c)
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.Matricula{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *MatriculaHandler) GetByID(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.Matricula{IDMatricula: id, IDEstudiante: 1, IDGrupo: 1, IDAnioLectivo: 1, Estado: "ACTIVO"})
}

func (h *MatriculaHandler) Create(c echo.Context) error {
	var m models.Matricula
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	m.IDMatricula = 1
	return c.JSON(http.StatusCreated, m)
}

func (h *MatriculaHandler) Update(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var m models.Matricula
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	m.IDMatricula = id
	return c.JSON(http.StatusOK, m)
}

func (h *MatriculaHandler) Delete(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	_ = id
	return c.NoContent(http.StatusNoContent)
}

type PeriodoHandler struct{}

func NewPeriodoHandler() *PeriodoHandler {
	return &PeriodoHandler{}
}

func (h *PeriodoHandler) GetAll(c echo.Context) error {
	offset, limit := parsePagination(c)
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.Periodo{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *PeriodoHandler) GetByID(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.Periodo{IDPeriodo: id, Nombre: "Primer Periodo"})
}

func (h *PeriodoHandler) Create(c echo.Context) error {
	var p models.Periodo
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	p.IDPeriodo = 1
	return c.JSON(http.StatusCreated, p)
}

func (h *PeriodoHandler) Update(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var p models.Periodo
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	p.IDPeriodo = id
	return c.JSON(http.StatusOK, p)
}

func (h *PeriodoHandler) Delete(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	_ = id
	return c.NoContent(http.StatusNoContent)
}

type AnioLectivoHandler struct{}

func NewAnioLectivoHandler() *AnioLectivoHandler {
	return &AnioLectivoHandler{}
}

func (h *AnioLectivoHandler) GetAll(c echo.Context) error {
	offset, limit := parsePagination(c)
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.AnioLectivo{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *AnioLectivoHandler) GetByID(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.AnioLectivo{IDAnioLectivo: id, Anio: 2025, Estado: "ACTIVO"})
}

func (h *AnioLectivoHandler) Create(c echo.Context) error {
	var a models.AnioLectivo
	if err := c.Bind(&a); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	a.IDAnioLectivo = 1
	return c.JSON(http.StatusCreated, a)
}

func (h *AnioLectivoHandler) Update(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var a models.AnioLectivo
	if err := c.Bind(&a); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	a.IDAnioLectivo = id
	return c.JSON(http.StatusOK, a)
}

func (h *AnioLectivoHandler) Delete(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	_ = id
	return c.NoContent(http.StatusNoContent)
}

type EvaluacionHandler struct{}

func NewEvaluacionHandler() *EvaluacionHandler {
	return &EvaluacionHandler{}
}

func (h *EvaluacionHandler) GetAll(c echo.Context) error {
	offset, limit := parsePagination(c)
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.Evaluacion{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *EvaluacionHandler) GetByID(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.Evaluacion{IDEvaluacion: id, Nombre: "Evaluación parcial"})
}

func (h *EvaluacionHandler) Create(c echo.Context) error {
	var e models.Evaluacion
	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	e.IDEvaluacion = 1
	return c.JSON(http.StatusCreated, e)
}

func (h *EvaluacionHandler) Update(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var e models.Evaluacion
	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	e.IDEvaluacion = id
	return c.JSON(http.StatusOK, e)
}

func (h *EvaluacionHandler) Delete(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	_ = id
	return c.NoContent(http.StatusNoContent)
}

type CalificacionHandler struct{}

func NewCalificacionHandler() *CalificacionHandler {
	return &CalificacionHandler{}
}

func (h *CalificacionHandler) GetAll(c echo.Context) error {
	offset, limit := parsePagination(c)
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.Calificacion{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *CalificacionHandler) GetByID(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.Calificacion{IDCalificacion: id, Nota: 4.5})
}

func (h *CalificacionHandler) Create(c echo.Context) error {
	var cals models.Calificacion
	if err := c.Bind(&cals); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	cals.IDCalificacion = 1
	return c.JSON(http.StatusCreated, cals)
}

func (h *CalificacionHandler) Update(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var cals models.Calificacion
	if err := c.Bind(&cals); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	cals.IDCalificacion = id
	return c.JSON(http.StatusOK, cals)
}

func (h *CalificacionHandler) Delete(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	_ = id
	return c.NoContent(http.StatusNoContent)
}

type CargaAcademicaHandler struct{}

func NewCargaAcademicaHandler() *CargaAcademicaHandler {
	return &CargaAcademicaHandler{}
}

func (h *CargaAcademicaHandler) GetAll(c echo.Context) error {
	offset, limit := parsePagination(c)
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.CargaAcademica{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *CargaAcademicaHandler) GetByID(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.CargaAcademica{IDCarga: id})
}

func (h *CargaAcademicaHandler) Create(c echo.Context) error {
	var ca models.CargaAcademica
	if err := c.Bind(&ca); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	ca.IDCarga = 1
	return c.JSON(http.StatusCreated, ca)
}

func (h *CargaAcademicaHandler) Update(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	var ca models.CargaAcademica
	if err := c.Bind(&ca); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Datos inválidos",
		})
	}
	ca.IDCarga = id
	return c.JSON(http.StatusOK, ca)
}

func (h *CargaAcademicaHandler) Delete(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	_ = id
	return c.NoContent(http.StatusNoContent)
}

type HorarioHandler struct{}

func NewHorarioHandler() *HorarioHandler {
	return &HorarioHandler{}
}

func (h *HorarioHandler) GetByDocente(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	_ = id
	return c.JSON(http.StatusOK, []models.Horario{})
}

func (h *HorarioHandler) GetByGrupo(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	_ = id
	return c.JSON(http.StatusOK, []models.Horario{})
}

type UsuarioHandler struct {
	repo *repository.UsuarioRepository
}

func NewUsuarioHandler(repo *repository.UsuarioRepository) *UsuarioHandler {
	return &UsuarioHandler{repo: repo}
}

func (h *UsuarioHandler) GetAll(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}

	usuarios, total, err := h.repo.GetAll(offset, limit)
	if err != nil {
		c.Logger().Errorf("Error fetching usuarios: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener los usuarios",
		})
	}

	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       usuarios,
		Pagination: models.Pagination{Total: total, Offset: offset, Limit: limit},
	})
}

func (h *UsuarioHandler) GetByID(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	usuario, err := h.repo.GetByID(id)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Usuario no encontrado",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al obtener el usuario",
		})
	}

	return c.JSON(http.StatusOK, usuario)
}

func (h *UsuarioHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}

	var req models.UsuarioCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	hash, err := utils.HashPassword(req.Clave)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al procesar la contraseña",
		})
	}

	u := models.Usuario{
		NombreUsuario: req.NombreUsuario,
		Clave:         hash,
		IDDocente:     req.IDDocente,
		IDEstudiante:  req.IDEstudiante,
		Activo:        req.Activo,
		IDRol:         req.IDRol,
	}

	existing, _ := h.repo.GetByUsername(req.NombreUsuario)
	if existing != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "El nombre de usuario ya existe",
		})
	}

	if err := h.repo.Create(&u); err != nil {
		c.Logger().Errorf("Error creating usuario: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al crear el usuario",
		})
	}

	claims := middleware.GetClaims(c)
	if claims != nil {
		c.Logger().Infof("AUDIT: Usuario '%s' (ID: %d) creó usuario '%s' (ID: %d)",
			claims.NombreUsuario, claims.IDUsuario, u.NombreUsuario, u.IDUsuario)
	}

	return c.JSON(http.StatusCreated, u.ToResponse())
}

func (h *UsuarioHandler) Update(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	var req models.UsuarioUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	u := models.Usuario{
		IDUsuario:     id,
		NombreUsuario: req.NombreUsuario,
		IDDocente:     req.IDDocente,
		IDEstudiante:  req.IDEstudiante,
		Activo:        req.Activo,
		IDRol:         req.IDRol,
	}

	if err := h.repo.Update(&u); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar the usuario",
		})
	}

	claims := middleware.GetClaims(c)
	if claims != nil {
		c.Logger().Infof("AUDIT: Usuario '%s' (ID: %d) actualizó usuario ID: %d",
			claims.NombreUsuario, claims.IDUsuario, id)
	}

	return c.JSON(http.StatusOK, u.ToResponse())
}

func (h *UsuarioHandler) Delete(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	claims := middleware.GetClaims(c)
	if claims != nil && claims.IDUsuario == id {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "No puedes eliminar tu propio usuario",
		})
	}

	if err := h.repo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al eliminar el usuario",
		})
	}

	if claims != nil {
		c.Logger().Infof("AUDIT: Usuario '%s' (ID: %d) eliminó usuario ID: %d",
			claims.NombreUsuario, claims.IDUsuario, id)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UsuarioHandler) ChangePassword(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}

	claims := middleware.GetClaims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "auth_error",
			Message: "Usuario no autenticado",
		})
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	isOwnPassword := claims.IDUsuario == targetID
	hasPermiso := false
	for _, p := range claims.Permisos {
		if p == "usuarios_cambiar_clave" {
			hasPermiso = true
			break
		}
	}

	if !hasPermiso && !isOwnPassword {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "No tienes permisos para cambiar esta contraseña",
		})
	}

	var req struct {
		Password string `json:"password" validate:"required,min=8,max=64"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al procesar la contraseña",
		})
	}

	if err := h.repo.UpdatePassword(targetID, hash); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar la contraseña",
		})
	}

	c.Logger().Infof("AUDIT: Usuario '%s' (ID: %d) cambió contraseña del usuario ID: %d",
		claims.NombreUsuario, claims.IDUsuario, targetID)

	return c.JSON(http.StatusOK, map[string]string{"message": "Password actualizado"})
}

type AcudienteHandler struct{}

func NewAcudienteHandler() *AcudienteHandler {
	return &AcudienteHandler{}
}

func (h *AcudienteHandler) GetAll(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.Acudiente{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *AcudienteHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, models.Acudiente{IDAcudiente: id})
}

func (h *AcudienteHandler) Create(c echo.Context) error {
	var a models.Acudiente
	c.Bind(&a)
	a.IDAcudiente = 1
	return c.JSON(http.StatusCreated, a)
}

func (h *AcudienteHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var a models.Acudiente
	c.Bind(&a)
	a.IDAcudiente = id
	return c.JSON(http.StatusOK, a)
}

func (h *AcudienteHandler) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
