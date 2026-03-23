package models

import "time"

type Institucion struct {
	IDInstitucion int    `json:"id_institucion" db:"id_institucion"`
	Nombre        string `json:"nombre" db:"nombre" validate:"required,min=5,max=100"`
	CodigoDane    string `json:"codigo_dane" db:"codigo_dane" validate:"required,regex=^[0-9]{12}$"`
}

type Sede struct {
	IDSede        int    `json:"id_sede" db:"id_sede"`
	Nombre        string `json:"nombre" db:"nombre" validate:"required,max=100"`
	Direccion     string `json:"direccion" db:"direccion" validate:"max=150"`
	IDInstitucion int    `json:"id_institucion" db:"id_institucion" validate:"required"`
	IDMunicipio   int    `json:"id_municipio" db:"id_municipio"`
}

type Grado struct {
	IDGrado int    `json:"id_grado" db:"id_grado"`
	Nombre  string `json:"nombre" db:"nombre" validate:"required,max=20"`
}

type Grupo struct {
	IDGrupo   int    `json:"id_grupo" db:"id_grupo"`
	Nombre    string `json:"nombre" db:"nombre" validate:"required,max=10"`
	IDGrado   int    `json:"id_grado" db:"id_grado" validate:"required"`
	IDSede    int    `json:"id_sede" db:"id_sede" validate:"required"`
	IDJornada int    `json:"id_jornada" db:"id_jornada" validate:"required"`
}

type Jornada struct {
	IDJornada int    `json:"id_jornada" db:"id_jornada"`
	Nombre    string `json:"nombre" db:"nombre" validate:"required,max=50"`
}

type Asignatura struct {
	IDAsignatura      int    `json:"id_asignatura" db:"id_asignatura"`
	Nombre            string `json:"nombre" db:"nombre" validate:"required,max=100"`
	IntensidadHoraria int    `json:"intensidad_horaria" db:"intensidad_horaria"`
}

type GradoAsignatura struct {
	IDGrado      int `json:"id_grado" db:"id_grado" validate:"required"`
	IDAsignatura int `json:"id_asignatura" db:"id_asignatura" validate:"required"`
}

type Estudiante struct {
	IDEstudiante       int        `json:"id_estudiante" db:"id_estudiante"`
	DocumentoIdentidad string     `json:"documento_identidad" db:"documento_identidad" validate:"required,regex=^[0-9]{8,10}$"`
	TipoDocumento      string     `json:"tipo_documento" db:"tipo_documento" validate:"required,oneof=CC TI RC CE PA"`
	FechaNacimiento    *time.Time `json:"fecha_nacimiento" db:"fecha_nacimiento"`
	Telefono           string     `json:"telefono" db:"telefono" validate:"regex=^\\+?[0-9]{7,15}$"`
	CorreoElectronico  string     `json:"correo_electronico" db:"correo_electronico" validate:"email,max=100"`
	Nombres            string     `json:"nombres" db:"nombres" validate:"required,max=100"`
	Apellidos          string     `json:"apellidos" db:"apellidos" validate:"required,max=100"`
	Direccion          string     `json:"direccion" db:"direccion" validate:"max=150"`
	IDMunicipio        int        `json:"id_municipio" db:"id_municipio"`
}

type Docente struct {
	IDDocente          int    `json:"id_docente" db:"id_docente"`
	Nombres            string `json:"nombres" db:"nombres" validate:"required,max=100"`
	Apellidos          string `json:"apellidos" db:"apellidos" validate:"required,max=100"`
	DocumentoIdentidad string `json:"documento_identidad" db:"documento_identidad" validate:"required,max=20"`
	TipoDocumento      string `json:"tipo_documento" db:"tipo_documento" validate:"required,oneof=CC TI CE PA"`
	Profesion          string `json:"profesion" db:"profesion" validate:"max=100"`
	Telefono           string `json:"telefono" db:"telefono" validate:"max=20"`
	CorreoElectronico  string `json:"correo_electronico" db:"correo_electronico" validate:"email,max=100"`
}

type CargaAcademica struct {
	IDCarga       int `json:"id_carga" db:"id_carga"`
	IDDocente     int `json:"id_docente" db:"id_docente" validate:"required"`
	IDGrupo       int `json:"id_grupo" db:"id_grupo" validate:"required"`
	IDAsignatura  int `json:"id_asignatura" db:"id_asignatura" validate:"required"`
	IDAnioLectivo int `json:"id_anio_lectivo" db:"id_anio_lectivo" validate:"required"`
}

type Periodo struct {
	IDPeriodo     int        `json:"id_periodo" db:"id_periodo"`
	Nombre        string     `json:"nombre" db:"nombre" validate:"required,max=50"`
	FechaInicio   *time.Time `json:"fecha_inicio" db:"fecha_inicio" validate:"required"`
	FechaFin      *time.Time `json:"fecha_fin" db:"fecha_fin" validate:"required"`
	IDAnioLectivo int        `json:"id_anio_lectivo" db:"id_anio_lectivo" validate:"required"`
}

type AnioLectivo struct {
	IDAnioLectivo int        `json:"id_anio_lectivo" db:"id_anio_lectivo"`
	Anio          int        `json:"anio" db:"anio" validate:"required"`
	FechaInicio   *time.Time `json:"fecha_inicio" db:"fecha_inicio" validate:"required"`
	FechaFin      *time.Time `json:"fecha_fin" db:"fecha_fin" validate:"required"`
	Estado        string     `json:"estado" db:"estado" validate:"oneof=ACTIVO CERRADO PLANEACION"`
}

type TipoEvaluacion struct {
	IDTipoEvaluacion int    `json:"id_tipo_evaluacion" db:"id_tipo_evaluacion"`
	Nombre           string `json:"nombre" db:"nombre" validate:"required,max=50"`
	Descripcion      string `json:"descripcion" db:"descripcion" validate:"max=255"`
}

type Evaluacion struct {
	IDEvaluacion      int        `json:"id_evaluacion" db:"id_evaluacion"`
	IDCargaAcademica  int        `json:"id_carga_academica" db:"id_carga_academica" validate:"required"`
	IDPeriodo         int        `json:"id_periodo" db:"id_periodo" validate:"required"`
	IDTipoEvaluacion  int        `json:"id_tipo_evaluacion" db:"id_tipo_evaluacion" validate:"required"`
	IDLogro           int        `json:"id_logro" db:"id_logro"`
	Nombre            string     `json:"nombre" db:"nombre" validate:"required,max=100"`
	Descripcion       string     `json:"descripcion" db:"descripcion"`
	FechaPresentacion *time.Time `json:"fecha_presentacion" db:"fecha_presentacion"`
}

type Pregunta struct {
	IDPregunta   int     `json:"id_pregunta" db:"id_pregunta"`
	IDEvaluacion int     `json:"id_evaluacion" db:"id_evaluacion" validate:"required"`
	Enunciado    string  `json:"enunciado" db:"enunciado" validate:"required,max=500"`
	Tipo         string  `json:"tipo" db:"tipo" validate:"required,oneof=ABIERTA OPCION_MULTIPLE VERDADERO_FALSO"`
	Peso         float64 `json:"peso" db:"peso" validate:"min=0.01,max=5.0"`
}

type OpcionRespuesta struct {
	IDOpcion   int    `json:"id_opcion" db:"id_opcion"`
	IDPregunta int    `json:"id_pregunta" db:"id_pregunta" validate:"required"`
	Texto      string `json:"texto" db:"texto" validate:"required,max=300"`
	EsCorrecta bool   `json:"es_correcta" db:"es_correcta"`
}

type Calificacion struct {
	IDCalificacion    int        `json:"id_calificacion" db:"id_calificacion"`
	IDEvaluacion      int        `json:"id_evaluacion" db:"id_evaluacion" validate:"required"`
	IDMatricula       int        `json:"id_matricula" db:"id_matricula" validate:"required"`
	Nota              float64    `json:"nota" db:"nota" validate:"min=0.0,max=5.0"`
	Observaciones     string     `json:"observaciones" db:"observaciones" validate:"max=255"`
	FechaCalificacion *time.Time `json:"fecha_calificacion" db:"fecha_calificacion"`
}

type RespuestaEstudiante struct {
	IDRespuesta          int     `json:"id_respuesta" db:"id_respuesta"`
	IDPregunta           int     `json:"id_pregunta" db:"id_pregunta" validate:"required"`
	IDCalificacion       int     `json:"id_calificacion" db:"id_calificacion" validate:"required"`
	Respuesta            string  `json:"respuesta" db:"respuesta" validate:"max=1000"`
	PuntajeObtenido      float64 `json:"puntaje_obtenido" db:"puntaje_obtenido"`
	IDOpcionSeleccionada int     `json:"id_opcion_seleccionada" db:"id_opcion_seleccionada"`
}

type Matricula struct {
	IDMatricula    int        `json:"id_matricula" db:"id_matricula"`
	IDEstudiante   int        `json:"id_estudiante" db:"id_estudiante" validate:"required"`
	IDGrupo        int        `json:"id_grupo" db:"id_grupo" validate:"required"`
	IDAnioLectivo  int        `json:"id_anio_lectivo" db:"id_anio_lectivo" validate:"required"`
	FechaMatricula *time.Time `json:"fecha_matricula" db:"fecha_matricula"`
	Estado         string     `json:"estado" db:"estado" validate:"oneof=ACTIVO DESERTOR PROMOVIDO RETIRADO GRADUADO"`
}

type Asistencia struct {
	IDAsistencia  int        `json:"id_asistencia" db:"id_asistencia"`
	IDMatricula   int        `json:"id_matricula" db:"id_matricula" validate:"required"`
	Fecha         *time.Time `json:"fecha" db:"fecha" validate:"required"`
	Estado        string     `json:"estado" db:"estado" validate:"required,oneof=PRESENTE AUSENTE TARDE JUSTIFICADO"`
	Observaciones string     `json:"observaciones" db:"observaciones" validate:"max=255"`
	FechaRegistro *time.Time `json:"fecha_registro" db:"fecha_registro"`
}

type ObservacionComportamiento struct {
	IDObservacion       int        `json:"id_observacion" db:"id_observacion"`
	IDMatricula         int        `json:"id_matricula" db:"id_matricula" validate:"required"`
	IDDocenteObservador int        `json:"id_docente_observador" db:"id_docente_observador" validate:"required"`
	IDCargaAcademica    int        `json:"id_carga_academica" db:"id_carga_academica"`
	FechaHora           *time.Time `json:"fecha_hora" db:"fecha_hora" validate:"required"`
	Descripcion         string     `json:"descripcion" db:"descripcion" validate:"required"`
	TipoObservacion     string     `json:"tipo_observacion" db:"tipo_observacion" validate:"max=50"`
	AccionesTomadas     string     `json:"acciones_tomadas" db:"acciones_tomadas" validate:"max=500"`
}

type Acudiente struct {
	IDAcudiente        int    `json:"id_acudiente" db:"id_acudiente"`
	Nombres            string `json:"nombres" db:"nombres" validate:"required,max=100"`
	Apellidos          string `json:"apellidos" db:"apellidos" validate:"required,max=100"`
	DocumentoIdentidad string `json:"documento_identidad" db:"documento_identidad" validate:"required,max=20"`
	TipoDocumento      string `json:"tipo_documento" db:"tipo_documento" validate:"required,oneof=CC CE PA"`
	Telefono           string `json:"telefono" db:"telefono" validate:"required,max=20"`
	Direccion          string `json:"direccion" db:"direccion" validate:"max=150"`
	CorreoElectronico  string `json:"correo_electronico" db:"correo_electronico" validate:"email,max=100"`
}

type EstudianteAcudiente struct {
	IDEstudiante int    `json:"id_estudiante" db:"id_estudiante" validate:"required"`
	IDAcudiente  int    `json:"id_acudiente" db:"id_acudiente" validate:"required"`
	Parentesco   string `json:"parentesco" db:"parentesco" validate:"max=50"`
	EsPrincipal  bool   `json:"es_principal" db:"es_principal"`
}

type Competencia struct {
	IDCompetencia int    `json:"id_competencia" db:"id_competencia"`
	Nombre        string `json:"nombre" db:"nombre" validate:"required,max=200"`
	Descripcion   string `json:"descripcion" db:"descripcion"`
	IDAsignatura  int    `json:"id_asignatura" db:"id_asignatura" validate:"required"`
	IDGrado       int    `json:"id_grado" db:"id_grado" validate:"required"`
}

type Logro struct {
	IDLogro       int     `json:"id_logro" db:"id_logro"`
	Descripcion   string  `json:"descripcion" db:"descripcion" validate:"required"`
	IDCompetencia int     `json:"id_competencia" db:"id_competencia" validate:"required"`
	IDPeriodo     int     `json:"id_periodo" db:"id_periodo" validate:"required"`
	Porcentaje    float64 `json:"porcentaje" db:"porcentaje" validate:"min=0,max=100"`
}

type ActividadRecuperacion struct {
	IDRecuperacion    int        `json:"id_recuperacion" db:"id_recuperacion"`
	IDCalificacion    int        `json:"id_calificacion" db:"id_calificacion" validate:"required"`
	Fecha             *time.Time `json:"fecha" db:"fecha" validate:"required"`
	Descripcion       string     `json:"descripcion" db:"descripcion" validate:"required"`
	NotaRecuperacion  float64    `json:"nota_recuperacion" db:"nota_recuperacion" validate:"min=0.0,max=5.0"`
	FechaCalificacion *time.Time `json:"fecha_calificacion" db:"fecha_calificacion"`
}

type Usuario struct {
	IDUsuario     int    `json:"id_usuario" db:"id_usuario"`
	NombreUsuario string `json:"nombre_usuario" db:"nombre_usuario" validate:"required,min=5,max=50,regex=^[a-zA-Z0-9_]+$"`
	Clave         string `json:"-" db:"clave" validate:"required,min=8,max=64"`
	IDDocente     *int   `json:"id_docente" db:"id_docente"`
	IDEstudiante  *int   `json:"id_estudiante" db:"id_estudiante"`
	Activo        bool   `json:"activo" db:"activo"`
	IDRol         int    `json:"id_rol" db:"id_rol" validate:"required"`
}

type UsuarioResponse struct {
	IDUsuario     int    `json:"id_usuario"`
	NombreUsuario string `json:"nombre_usuario"`
	IDDocente     *int   `json:"id_docente,omitempty"`
	IDEstudiante  *int   `json:"id_estudiante,omitempty"`
	Activo        bool   `json:"activo"`
	IDRol         int    `json:"id_rol"`
}

type UsuarioCreateRequest struct {
	NombreUsuario string `json:"nombre_usuario" validate:"required,min=5,max=50"`
	Clave         string `json:"clave" validate:"required,min=8,max=64"`
	IDDocente     *int   `json:"id_docente,omitempty"`
	IDEstudiante  *int   `json:"id_estudiante,omitempty"`
	Activo        bool   `json:"activo"`
	IDRol         int    `json:"id_rol" validate:"required"`
}

type UsuarioUpdateRequest struct {
	NombreUsuario string `json:"nombre_usuario" validate:"required,min=5,max=50"`
	IDDocente     *int   `json:"id_docente,omitempty"`
	IDEstudiante  *int   `json:"id_estudiante,omitempty"`
	Activo        bool   `json:"activo"`
	IDRol         int    `json:"id_rol" validate:"required"`
}

func (u *Usuario) ToResponse() UsuarioResponse {
	return UsuarioResponse{
		IDUsuario:     u.IDUsuario,
		NombreUsuario: u.NombreUsuario,
		IDDocente:     u.IDDocente,
		IDEstudiante:  u.IDEstudiante,
		Activo:        u.Activo,
		IDRol:         u.IDRol,
	}
}

func UsuariosToResponse(usuarios []Usuario) []UsuarioResponse {
	result := make([]UsuarioResponse, len(usuarios))
	for i, u := range usuarios {
		result[i] = u.ToResponse()
	}
	return result
}

type Rol struct {
	IDRol         int       `json:"id_rol" db:"id_rol"`
	Nombre        string    `json:"nombre" db:"nombre" validate:"required,max=50"`
	Descripcion   string    `json:"descripcion" db:"descripcion" validate:"max=255"`
	EsRolSistema  bool      `json:"es_rol_sistema" db:"es_rol_sistema"`
	PermisosCount int       `json:"permisos_count,omitempty"`
	Permisos      []Permiso `json:"permisos,omitempty"`
}

type Modulo struct {
	IDModulo      int    `json:"id_modulo" db:"id_modulo"`
	Nombre        string `json:"nombre" db:"nombre" validate:"required,max=50"`
	Descripcion   string `json:"descripcion" db:"descripcion" validate:"max=255"`
	Codigo        string `json:"codigo" db:"codigo" validate:"required,max=30"`
	PermisosCount int    `json:"permisos_count,omitempty"`
}

type Permiso struct {
	IDPermiso    int    `json:"id_permiso" db:"id_permiso"`
	Nombre       string `json:"nombre" db:"nombre" validate:"required,max=100"`
	Descripcion  string `json:"descripcion" db:"descripcion" validate:"max=255"`
	Codigo       string `json:"codigo" db:"codigo" validate:"required,max=30"`
	IDModulo     int    `json:"id_modulo" db:"id_modulo" validate:"required"`
	ModuloNombre string `json:"modulo_nombre,omitempty"`
	ModuloCodigo string `json:"modulo_codigo,omitempty"`
}

type RolPermiso struct {
	IDRol     int `json:"id_rol" db:"id_rol" validate:"required"`
	IDPermiso int `json:"id_permiso" db:"id_permiso" validate:"required"`
}

type Departamento struct {
	IDDepartamento int    `json:"id_departamento" db:"id_departamento"`
	Nombre         string `json:"nombre" db:"nombre" validate:"required,max=100"`
	Codigo         string `json:"codigo" db:"codigo" validate:"required,max=5"`
}

type Municipio struct {
	IDMunicipio    int    `json:"id_municipio" db:"id_municipio"`
	Nombre         string `json:"nombre" db:"nombre" validate:"required,max=100"`
	Codigo         string `json:"codigo" db:"codigo" validate:"required,max=5"`
	IDDepartamento int    `json:"id_departamento" db:"id_departamento" validate:"required"`
}

type DiaSemana struct {
	IDDia  int    `json:"id_dia" db:"id_dia"`
	Nombre string `json:"nombre" db:"nombre" validate:"required,max=20"`
	Codigo int    `json:"codigo" db:"codigo" validate:"min=1,max=7"`
}

type Horario struct {
	IDHorario        int    `json:"id_horario" db:"id_horario"`
	IDCargaAcademica int    `json:"id_carga_academica" db:"id_carga_academica" validate:"required"`
	IDDia            int    `json:"id_dia" db:"id_dia" validate:"required,min=1,max=7"`
	HoraInicio       string `json:"hora_inicio" db:"hora_inicio" validate:"required"`
	HoraFin          string `json:"hora_fin" db:"hora_fin" validate:"required"`
}

type ArchivoDigital struct {
	IDArchivo          int        `json:"id_archivo" db:"id_archivo"`
	TipoArchivo        string     `json:"tipo_archivo" db:"tipo_archivo" validate:"required,max=50"`
	NombreArchivo      string     `json:"nombre_archivo" db:"nombre_archivo" validate:"required,max=255"`
	RutaAlmacenamiento string     `json:"ruta_almacenamiento" db:"ruta_almacenamiento" validate:"required,max=500"`
	MimeType           string     `json:"mime_type" db:"mime_type" validate:"required,max=100"`
	FechaCarga         *time.Time `json:"fecha_carga" db:"fecha_carga"`
	IDUsuarioCarga     int        `json:"id_usuario_carga" db:"id_usuario_carga"`
	EntidadRelacionada string     `json:"entidad_relacionada" db:"entidad_relacionada" validate:"required,max=50"`
	IDEntidad          int        `json:"id_entidad" db:"id_entidad" validate:"required"`
}

type ErrorResponse struct {
	Error   string       `json:"error"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type Pagination struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int      `json:"expires_in"`
	NombreUsuario string  `json:"nombre_usuario"`
	Role         string   `json:"role"`
	IDRol        int      `json:"id_rol"`
	Permisos     []string `json:"permisos"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"` // Opcional, ahora preferimos cookies
}

type RefreshTokenResponse struct {
	AccessToken  string   `json:"access_token"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int      `json:"expires_in"`
	NombreUsuario string  `json:"nombre_usuario"`
	Role         string   `json:"role"`
	IDRol        int      `json:"id_rol"`
	Permisos     []string `json:"permisos"`
}

type RefreshTokenDB struct {
	ID              int       `db:"id"`
	TokenHash       string    `db:"token_hash"`
	JTI             string    `db:"jti"`
	IDUsuario       int       `db:"id_usuario"`
	FechaExpiracion time.Time `db:"fecha_expiracion"`
	FechaCreacion   time.Time `db:"fecha_creacion"`
	Dispositivo     *string   `db:"dispositivo"`
	Activo          bool      `db:"activo"`
}

type RevokedToken struct {
	JTI             string    `db:"jti"`
	FechaRevocado   time.Time `db:"fecha_revocado"`
	FechaExpiracion time.Time `db:"fecha_expiracion"`
}
