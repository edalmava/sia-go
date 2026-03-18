package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/edalmava/sia/internal/config"
	"github.com/edalmava/sia/internal/database"
	"github.com/edalmava/sia/internal/handlers"
	"github.com/edalmava/sia/internal/middleware"
	"github.com/edalmava/sia/internal/repository"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	e := echo.New()

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/web/login.html")
	})

	e.Static("/web", "web")

	e.GET("/web/login.html", func(c echo.Context) error {
		return c.File(filepath.Join("web", "login.html"))
	})

	e.GET("/web/dashboard.html", func(c echo.Context) error {
		return c.File(filepath.Join("web", "dashboard.html"))
	})

	db, err := database.Connect(&database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Printf("Warning: Could not connect to database: %v", err)
		log.Println("Running without database connection")
	} else {
		defer db.Close()

		if err := db.Migrate(); err != nil {
			log.Printf("Warning: Migration failed: %v", err)
		} else {
			log.Println("Database migrations completed successfully")
		}
	}

	var repo *repository.Repository
	if db != nil {
		repo = repository.NewRepository(db.DB)
	} else {
		repo = repository.NewRepository(nil)
	}

	instHandler := handlers.NewInstitucionHandler(repo.Institucion)
	sedeHandler := handlers.NewSedeHandler(repo.Sede)
	gradoHandler := handlers.NewGradoHandler(repo.Grado)
	grupoHandler := handlers.NewGrupoHandler(repo.Grupo)
	jornadaHandler := handlers.NewJornadaHandler(repo.Jornada)
	asignaturaHandler := handlers.NewAsignaturaHandler(repo.Asignatura)
	estudianteHandler := handlers.NewEstudianteHandler(repo.Estudiante)
	docenteHandler := handlers.NewDocenteHandler(repo.Docente)
	matriculaHandler := handlers.NewMatriculaHandler()
	periodoHandler := handlers.NewPeriodoHandler()
	anioHandler := handlers.NewAnioLectivoHandler()
	evaluacionHandler := handlers.NewEvaluacionHandler()
	calificacionHandler := handlers.NewCalificacionHandler()
	cargaHandler := handlers.NewCargaAcademicaHandler()
	horarioHandler := handlers.NewHorarioHandler()
	usuarioHandler := handlers.NewUsuarioHandler(repo.Usuario)
	acudienteHandler := handlers.NewAcudienteHandler()
	authHandler := handlers.NewAuthHandler(cfg, repo.Usuario, repo.Permiso)
	configHandler := handlers.NewConfigHandler(repo.Rol, repo.Permiso, repo.Modulo)

	e.POST("/auth/login", authHandler.Login)
	e.POST("/auth/refresh", authHandler.Refresh)
	e.POST("/auth/logout", authHandler.Logout)

	api := e.Group("/api/v1")
	api.Use(middleware.JWTAuth(cfg))

	adminAPI := api.Group("")
	adminAPI.Use(middleware.RequireRole("ADMIN"))

	api.GET("/instituciones", instHandler.GetAll)
	api.GET("/instituciones/:id", instHandler.GetByID)
	api.POST("/instituciones", instHandler.Create)
	api.PUT("/instituciones/:id", instHandler.Update)
	api.DELETE("/instituciones/:id", instHandler.Delete)

	api.GET("/sedes", sedeHandler.GetAll)
	api.GET("/sedes/:id", sedeHandler.GetByID)
	api.POST("/sedes", sedeHandler.Create)
	api.PUT("/sedes/:id", sedeHandler.Update)
	api.DELETE("/sedes/:id", sedeHandler.Delete)

	api.GET("/grados", gradoHandler.GetAll)
	api.GET("/grados/:id", gradoHandler.GetByID)
	api.POST("/grados", gradoHandler.Create)
	api.PUT("/grados/:id", gradoHandler.Update)
	api.DELETE("/grados/:id", gradoHandler.Delete)
	api.GET("/grados/:id/asignaturas", gradoHandler.GetAsignaturas)
	api.POST("/grades/:id/asignaturas", gradoHandler.AddAsignatura)
	api.DELETE("/grades/:id_grado/asignaturas/:id_asignatura", gradoHandler.RemoveAsignatura)

	api.GET("/grupos", grupoHandler.GetAll)
	api.GET("/grupos/:id", grupoHandler.GetByID)
	api.POST("/grupos", grupoHandler.Create)
	api.PUT("/grupos/:id", grupoHandler.Update)
	api.DELETE("/grupos/:id", grupoHandler.Delete)

	api.GET("/jornadas", jornadaHandler.GetAll)
	api.GET("/jornadas/:id", jornadaHandler.GetByID)
	api.POST("/jornadas", jornadaHandler.Create)
	api.PUT("/jornadas/:id", jornadaHandler.Update)
	api.DELETE("/jornadas/:id", jornadaHandler.Delete)

	api.GET("/asignaturas", asignaturaHandler.GetAll)
	api.GET("/asignaturas/:id", asignaturaHandler.GetByID)
	api.POST("/asignaturas", asignaturaHandler.Create)
	api.PUT("/asignaturas/:id", asignaturaHandler.Update)
	api.DELETE("/asignaturas/:id", asignaturaHandler.Delete)

	api.GET("/estudiantes", estudianteHandler.GetAll)
	api.GET("/estudiantes/:id", estudianteHandler.GetByID)
	api.POST("/estudiantes", estudianteHandler.Create)
	api.PUT("/estudiantes/:id", estudianteHandler.Update)
	api.DELETE("/estudiantes/:id", estudianteHandler.Delete)

	api.GET("/docentes", docenteHandler.GetAll)
	api.GET("/docentes/:id", docenteHandler.GetByID)
	api.POST("/docentes", docenteHandler.Create)
	api.PUT("/docentes/:id", docenteHandler.Update)
	api.DELETE("/docentes/:id", docenteHandler.Delete)

	api.GET("/matriculas", matriculaHandler.GetAll)
	api.GET("/matriculas/:id", matriculaHandler.GetByID)
	api.POST("/matriculas", matriculaHandler.Create)
	api.PUT("/matriculas/:id", matriculaHandler.Update)
	api.DELETE("/matriculas/:id", matriculaHandler.Delete)

	api.GET("/periodos", periodoHandler.GetAll)
	api.GET("/periodos/:id", periodoHandler.GetByID)
	api.POST("/periodos", periodoHandler.Create)
	api.PUT("/periodos/:id", periodoHandler.Update)
	api.DELETE("/periodos/:id", periodoHandler.Delete)

	api.GET("/anios-lectivos", anioHandler.GetAll)
	api.GET("/anios-lectivos/:id", anioHandler.GetByID)
	api.POST("/anios-lectivos", anioHandler.Create)
	api.PUT("/anios-lectivos/:id", anioHandler.Update)
	api.DELETE("/anios-lectivos/:id", anioHandler.Delete)

	api.GET("/evaluaciones", evaluacionHandler.GetAll)
	api.GET("/evaluaciones/:id", evaluacionHandler.GetByID)
	api.POST("/evaluaciones", evaluacionHandler.Create)
	api.PUT("/evaluaciones/:id", evaluacionHandler.Update)
	api.DELETE("/evaluaciones/:id", evaluacionHandler.Delete)

	api.GET("/calificaciones", calificacionHandler.GetAll)
	api.GET("/calificaciones/:id", calificacionHandler.GetByID)
	api.POST("/calificaciones", calificacionHandler.Create)
	api.PUT("/calificaciones/:id", calificacionHandler.Update)
	api.DELETE("/calificaciones/:id", calificacionHandler.Delete)

	api.GET("/cargas-academicas", cargaHandler.GetAll)
	api.GET("/cargas-academicas/:id", cargaHandler.GetByID)
	api.POST("/cargas-academicas", cargaHandler.Create)
	api.PUT("/cargas-academicas/:id", cargaHandler.Update)
	api.DELETE("/cargas-academicas/:id", cargaHandler.Delete)

	api.GET("/docentes/:id/horarios", horarioHandler.GetByDocente)
	api.GET("/grupos/:id/horarios", horarioHandler.GetByGrupo)

	adminAPI.GET("/usuarios", usuarioHandler.GetAll)
	adminAPI.GET("/usuarios/:id", usuarioHandler.GetByID)
	adminAPI.POST("/usuarios", usuarioHandler.Create)
	adminAPI.PUT("/usuarios/:id", usuarioHandler.Update)
	adminAPI.DELETE("/usuarios/:id", usuarioHandler.Delete)
	adminAPI.POST("/usuarios/:id/change-password", usuarioHandler.ChangePassword)

	api.GET("/acudientes", acudienteHandler.GetAll)
	api.GET("/acudientes/:id", acudienteHandler.GetByID)
	api.POST("/acudientes", acudienteHandler.Create)
	api.PUT("/acudientes/:id", acudienteHandler.Update)
	api.DELETE("/acudientes/:id", acudienteHandler.Delete)

	adminAPI.GET("/roles", configHandler.GetRoles)
	adminAPI.GET("/roles/:id", configHandler.GetRoleByID)
	adminAPI.POST("/roles", configHandler.CreateRole)
	adminAPI.PUT("/roles/:id", configHandler.UpdateRole)
	adminAPI.DELETE("/roles/:id", configHandler.DeleteRole)

	adminAPI.GET("/permisos", configHandler.GetPermisos)
	adminAPI.GET("/modulos", configHandler.GetModulos)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Servidor iniciando en %s", addr)
	e.Logger.Fatal(e.Start(addr))
}
