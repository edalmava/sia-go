package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/edalmava/sia/internal/middleware"
	"github.com/edalmava/sia/internal/models"
	"github.com/edalmava/sia/internal/repository"
	"github.com/edalmava/sia/internal/utils"
	"github.com/labstack/echo/v4"
)

type MatriculaHandler struct{}

func NewMatriculaHandler() *MatriculaHandler {
	return &MatriculaHandler{}
}

func (h *MatriculaHandler) GetAll(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.Matricula{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *MatriculaHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, models.Matricula{IDMatricula: id, IDEstudiante: 1, IDGrupo: 1, IDAnioLectivo: 1, Estado: "ACTIVO"})
}

func (h *MatriculaHandler) Create(c echo.Context) error {
	var m models.Matricula
	c.Bind(&m)
	m.IDMatricula = 1
	return c.JSON(http.StatusCreated, m)
}

func (h *MatriculaHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var m models.Matricula
	c.Bind(&m)
	m.IDMatricula = id
	return c.JSON(http.StatusOK, m)
}

func (h *MatriculaHandler) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

type PeriodoHandler struct{}

func NewPeriodoHandler() *PeriodoHandler {
	return &PeriodoHandler{}
}

func (h *PeriodoHandler) GetAll(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.Periodo{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *PeriodoHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, models.Periodo{IDPeriodo: id, Nombre: "Primer Periodo"})
}

func (h *PeriodoHandler) Create(c echo.Context) error {
	var p models.Periodo
	c.Bind(&p)
	p.IDPeriodo = 1
	return c.JSON(http.StatusCreated, p)
}

func (h *PeriodoHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var p models.Periodo
	c.Bind(&p)
	p.IDPeriodo = id
	return c.JSON(http.StatusOK, p)
}

func (h *PeriodoHandler) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

type AnioLectivoHandler struct{}

func NewAnioLectivoHandler() *AnioLectivoHandler {
	return &AnioLectivoHandler{}
}

func (h *AnioLectivoHandler) GetAll(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.AnioLectivo{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *AnioLectivoHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, models.AnioLectivo{IDAnioLectivo: id, Anio: 2025, Estado: "ACTIVO"})
}

func (h *AnioLectivoHandler) Create(c echo.Context) error {
	var a models.AnioLectivo
	c.Bind(&a)
	a.IDAnioLectivo = 1
	return c.JSON(http.StatusCreated, a)
}

func (h *AnioLectivoHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var a models.AnioLectivo
	c.Bind(&a)
	a.IDAnioLectivo = id
	return c.JSON(http.StatusOK, a)
}

func (h *AnioLectivoHandler) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

type EvaluacionHandler struct{}

func NewEvaluacionHandler() *EvaluacionHandler {
	return &EvaluacionHandler{}
}

func (h *EvaluacionHandler) GetAll(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.Evaluacion{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *EvaluacionHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, models.Evaluacion{IDEvaluacion: id, Nombre: "Evaluación parcial"})
}

func (h *EvaluacionHandler) Create(c echo.Context) error {
	var e models.Evaluacion
	c.Bind(&e)
	e.IDEvaluacion = 1
	return c.JSON(http.StatusCreated, e)
}

func (h *EvaluacionHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var e models.Evaluacion
	c.Bind(&e)
	e.IDEvaluacion = id
	return c.JSON(http.StatusOK, e)
}

func (h *EvaluacionHandler) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

type CalificacionHandler struct{}

func NewCalificacionHandler() *CalificacionHandler {
	return &CalificacionHandler{}
}

func (h *CalificacionHandler) GetAll(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.Calificacion{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *CalificacionHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, models.Calificacion{IDCalificacion: id, Nota: 4.5})
}

func (h *CalificacionHandler) Create(c echo.Context) error {
	var cals models.Calificacion
	c.Bind(&cals)
	cals.IDCalificacion = 1
	return c.JSON(http.StatusCreated, cals)
}

func (h *CalificacionHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var cals models.Calificacion
	c.Bind(&cals)
	cals.IDCalificacion = id
	return c.JSON(http.StatusOK, cals)
}

func (h *CalificacionHandler) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

type CargaAcademicaHandler struct{}

func NewCargaAcademicaHandler() *CargaAcademicaHandler {
	return &CargaAcademicaHandler{}
}

func (h *CargaAcademicaHandler) GetAll(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       []models.CargaAcademica{},
		Pagination: models.Pagination{Total: 0, Offset: offset, Limit: limit},
	})
}

func (h *CargaAcademicaHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, models.CargaAcademica{IDCarga: id})
}

func (h *CargaAcademicaHandler) Create(c echo.Context) error {
	var ca models.CargaAcademica
	c.Bind(&ca)
	ca.IDCarga = 1
	return c.JSON(http.StatusCreated, ca)
}

func (h *CargaAcademicaHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var ca models.CargaAcademica
	c.Bind(&ca)
	ca.IDCarga = id
	return c.JSON(http.StatusOK, ca)
}

func (h *CargaAcademicaHandler) Delete(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

type HorarioHandler struct{}

func NewHorarioHandler() *HorarioHandler {
	return &HorarioHandler{}
}

func (h *HorarioHandler) GetByDocente(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	_ = id
	return c.JSON(http.StatusOK, []models.Horario{})
}

func (h *HorarioHandler) GetByGrupo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
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

func hasPermission(permisos []string, permiso string) bool {
	for _, p := range permisos {
		if p == permiso {
			return true
		}
	}
	return false
}

func checkPermission(c echo.Context, permiso string) bool {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return false
	}
	return hasPermission(claims.Permisos, permiso)
}

func (h *UsuarioHandler) Create(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}

	if !checkPermission(c, "usuarios_crear") {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "No tienes permisos para crear usuarios",
		})
	}

	var u models.Usuario
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	if u.Clave == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "La contraseña es requerida",
		})
	}

	hash, err := utils.HashPassword(u.Clave)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al procesar la contraseña",
		})
	}
	u.Clave = hash

	if err := h.repo.Create(&u); err != nil {
		c.Logger().Errorf("Error creating usuario: %v", err)

		errStr := err.Error()
		if strings.Contains(errStr, "usuarios_nombre_usuario_key") || strings.Contains(errStr, "duplicate key") {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "validation_error",
				Message: "El nombre de usuario ya existe",
			})
		}

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

	u.Clave = ""
	return c.JSON(http.StatusCreated, u)
}

func (h *UsuarioHandler) Update(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}

	if !checkPermission(c, "usuarios_editar") {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "No tienes permisos para actualizar usuarios",
		})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "ID inválido",
		})
	}

	var u models.Usuario
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Error en los datos de entrada",
		})
	}

	u.IDUsuario = id

	if err := h.repo.Update(&u); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Error al actualizar el usuario",
		})
	}

	claims := middleware.GetClaims(c)
	if claims != nil {
		c.Logger().Infof("AUDIT: Usuario '%s' (ID: %d) actualizó usuario ID: %d",
			claims.NombreUsuario, claims.IDUsuario, id)
	}

	return c.JSON(http.StatusOK, u)
}

func (h *UsuarioHandler) Delete(c echo.Context) error {
	if h.repo == nil {
		return dbUnavailable(c)
	}

	if !checkPermission(c, "usuarios_eliminar") {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "No tienes permisos para eliminar usuarios",
		})
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

	if claims := middleware.GetClaims(c); claims != nil {
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

	hasPermiso := checkPermission(c, "usuarios_cambiar_clave")
	isOwnPassword := claims.IDUsuario == targetID

	if !hasPermiso && !isOwnPassword {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error:   "forbidden",
			Message: "No tienes permisos para cambiar esta contraseña",
		})
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil || req.Password == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Contraseña requerida",
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
