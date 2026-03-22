package database

import "fmt"

func (db *DB) Migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS departamentos (
			id_departamento SERIAL PRIMARY KEY,
			nombre VARCHAR(100) NOT NULL,
			codigo VARCHAR(5) NOT NULL UNIQUE
		)`,

		`CREATE TABLE IF NOT EXISTS municipios (
			id_municipio SERIAL PRIMARY KEY,
			nombre VARCHAR(100) NOT NULL,
			codigo VARCHAR(5) NOT NULL,
			id_departamento INTEGER NOT NULL REFERENCES departamentos(id_departamento),
			UNIQUE(nombre, id_departamento)
		)`,

		`CREATE TABLE IF NOT EXISTS instituciones (
			id_institucion SERIAL PRIMARY KEY,
			nombre VARCHAR(100) NOT NULL,
			codigo_dane VARCHAR(12) NOT NULL UNIQUE
		)`,

		`CREATE TABLE IF NOT EXISTS sedes (
			id_sede SERIAL PRIMARY KEY,
			nombre VARCHAR(100) NOT NULL,
			direccion VARCHAR(150),
			id_institucion INTEGER NOT NULL REFERENCES instituciones(id_institucion),
			id_municipio INTEGER REFERENCES municipios(id_municipio)
		)`,

		`CREATE TABLE IF NOT EXISTS jornadas (
			id_jornada SERIAL PRIMARY KEY,
			nombre VARCHAR(50) NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS grados (
			id_grado SERIAL PRIMARY KEY,
			nombre VARCHAR(20) NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS grupos (
			id_grupo SERIAL PRIMARY KEY,
			nombre VARCHAR(10) NOT NULL,
			id_grado INTEGER NOT NULL REFERENCES grados(id_grado),
			id_sede INTEGER NOT NULL REFERENCES sedes(id_sede),
			id_jornada INTEGER NOT NULL REFERENCES jornadas(id_jornada),
			UNIQUE(nombre, id_grado, id_sede, id_jornada)
		)`,

		`CREATE TABLE IF NOT EXISTS asignaturas (
			id_asignatura SERIAL PRIMARY KEY,
			nombre VARCHAR(100) NOT NULL,
			intensidad_horaria INTEGER
		)`,

		`CREATE TABLE IF NOT EXISTS grados_asignaturas (
			id_grado INTEGER NOT NULL REFERENCES grados(id_grado),
			id_asignatura INTEGER NOT NULL REFERENCES asignaturas(id_asignatura),
			PRIMARY KEY (id_grado, id_asignatura)
		)`,

		`CREATE TABLE IF NOT EXISTS docentes (
			id_docente SERIAL PRIMARY KEY,
			nombres VARCHAR(100) NOT NULL,
			apellidos VARCHAR(100) NOT NULL,
			documento_identidad VARCHAR(20) NOT NULL UNIQUE,
			tipo_documento VARCHAR(2) NOT NULL,
			profesion VARCHAR(100),
			telefono VARCHAR(20),
			correo_electronico VARCHAR(100)
		)`,

		`CREATE TABLE IF NOT EXISTS estudiantes (
			id_estudiante SERIAL PRIMARY KEY,
			documento_identidad VARCHAR(10) NOT NULL UNIQUE,
			tipo_documento VARCHAR(2) NOT NULL,
			fecha_nacimiento DATE,
			telefono VARCHAR(15),
			correo_electronico VARCHAR(100),
			nombres VARCHAR(100) NOT NULL,
			apellidos VARCHAR(100) NOT NULL,
			direccion VARCHAR(150),
			id_municipio INTEGER REFERENCES municipios(id_municipio)
		)`,

		`CREATE TABLE IF NOT EXISTS anios_lectivos (
			id_anio_lectivo SERIAL PRIMARY KEY,
			anio INTEGER NOT NULL,
			fecha_inicio DATE NOT NULL,
			fecha_fin DATE NOT NULL,
			estado VARCHAR(20) DEFAULT 'PLANEACION',
			UNIQUE(anio)
		)`,

		`CREATE TABLE IF NOT EXISTS periodos (
			id_periodo SERIAL PRIMARY KEY,
			nombre VARCHAR(50) NOT NULL,
			fecha_inicio DATE NOT NULL,
			fecha_fin DATE NOT NULL,
			id_anio_lectivo INTEGER NOT NULL REFERENCES anios_lectivos(id_anio_lectivo)
		)`,

		`CREATE TABLE IF NOT EXISTS cargas_academicas (
			id_carga SERIAL PRIMARY KEY,
			id_docente INTEGER NOT NULL REFERENCES docentes(id_docente),
			id_grupo INTEGER NOT NULL REFERENCES grupos(id_grupo),
			id_asignatura INTEGER NOT NULL REFERENCES asignaturas(id_asignatura),
			id_anio_lectivo INTEGER NOT NULL REFERENCES anios_lectivos(id_anio_lectivo),
			UNIQUE(id_docente, id_grupo, id_asignatura, id_anio_lectivo)
		)`,

		`CREATE TABLE IF NOT EXISTS tipo_evaluaciones (
			id_tipo_evaluacion SERIAL PRIMARY KEY,
			nombre VARCHAR(50) NOT NULL,
			descripcion VARCHAR(255)
		)`,

		`CREATE TABLE IF NOT EXISTS competencias (
			id_competencia SERIAL PRIMARY KEY,
			nombre VARCHAR(200) NOT NULL,
			descripcion TEXT,
			id_asignatura INTEGER NOT NULL REFERENCES asignaturas(id_asignatura),
			id_grado INTEGER NOT NULL REFERENCES grados(id_grado)
		)`,

		`CREATE TABLE IF NOT EXISTS logros (
			id_logro SERIAL PRIMARY KEY,
			descripcion TEXT NOT NULL,
			id_competencia INTEGER NOT NULL REFERENCES competencias(id_competencia),
			id_periodo INTEGER NOT NULL REFERENCES periodos(id_periodo),
			porcentaje DECIMAL(5,2)
		)`,

		`CREATE TABLE IF NOT EXISTS evaluaciones (
			id_evaluacion SERIAL PRIMARY KEY,
			id_carga_academica INTEGER NOT NULL REFERENCES cargas_academicas(id_carga),
			id_periodo INTEGER NOT NULL REFERENCES periodos(id_periodo),
			id_tipo_evaluacion INTEGER NOT NULL REFERENCES tipo_evaluaciones(id_tipo_evaluacion),
			id_logro INTEGER REFERENCES logros(id_logro),
			nombre VARCHAR(100) NOT NULL,
			descripcion TEXT,
			fecha_presentacion DATE
		)`,

		`CREATE TABLE IF NOT EXISTS preguntas (
			id_pregunta SERIAL PRIMARY KEY,
			id_evaluacion INTEGER NOT NULL REFERENCES evaluaciones(id_evaluacion),
			enunciado VARCHAR(500) NOT NULL,
			tipo VARCHAR(20) NOT NULL,
			peso DECIMAL(4,2) DEFAULT 1.0
		)`,

		`CREATE TABLE IF NOT EXISTS opciones_respuestas (
			id_opcion SERIAL PRIMARY KEY,
			id_pregunta INTEGER NOT NULL REFERENCES preguntas(id_pregunta),
			texto VARCHAR(300) NOT NULL,
			es_correcta BOOLEAN DEFAULT FALSE
		)`,

		`CREATE TABLE IF NOT EXISTS matriculas (
			id_matricula SERIAL PRIMARY KEY,
			id_estudiante INTEGER NOT NULL REFERENCES estudiantes(id_estudiante),
			id_grupo INTEGER NOT NULL REFERENCES grupos(id_grupo),
			id_anio_lectivo INTEGER NOT NULL REFERENCES anios_lectivos(id_anio_lectivo),
			fecha_matricula DATE,
			estado VARCHAR(20) DEFAULT 'ACTIVO',
			UNIQUE(id_estudiante, id_grupo, id_anio_lectivo)
		)`,

		`CREATE TABLE IF NOT EXISTS calificaciones (
			id_calificacion SERIAL PRIMARY KEY,
			id_evaluacion INTEGER NOT NULL REFERENCES evaluaciones(id_evaluacion),
			id_matricula INTEGER NOT NULL REFERENCES matriculas(id_matricula),
			nota DECIMAL(3,2) NOT NULL,
			observaciones VARCHAR(255),
			fecha_calificacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS respuesta_estudiantes (
			id_respuesta SERIAL PRIMARY KEY,
			id_pregunta INTEGER NOT NULL REFERENCES preguntas(id_pregunta),
			id_calificacion INTEGER NOT NULL REFERENCES calificaciones(id_calificacion),
			respuesta VARCHAR(1000),
			puntaje_obtenido DECIMAL(4,2),
			id_opcion_seleccionada INTEGER REFERENCES opciones_respuestas(id_opcion)
		)`,

		`CREATE TABLE IF NOT EXISTS asistencia (
			id_asistencia SERIAL PRIMARY KEY,
			id_matricula INTEGER NOT NULL REFERENCES matriculas(id_matricula),
			fecha DATE NOT NULL,
			estado VARCHAR(20) NOT NULL,
			observaciones VARCHAR(255),
			fecha_registro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(id_matricula, fecha)
		)`,

		`CREATE TABLE IF NOT EXISTS observacion_comportamiento (
			id_observacion SERIAL PRIMARY KEY,
			id_matricula INTEGER NOT NULL REFERENCES matriculas(id_matricula),
			id_docente_observador INTEGER NOT NULL REFERENCES docentes(id_docente),
			id_carga_academica INTEGER REFERENCES cargas_academicas(id_carga),
			fecha_hora TIMESTAMP NOT NULL,
			descripcion TEXT NOT NULL,
			tipo_observacion VARCHAR(50),
			acciones_tomadas VARCHAR(500)
		)`,

		`CREATE TABLE IF NOT EXISTS acudientes (
			id_acudiente SERIAL PRIMARY KEY,
			nombres VARCHAR(100) NOT NULL,
			apellidos VARCHAR(100) NOT NULL,
			documento_identidad VARCHAR(20) NOT NULL UNIQUE,
			tipo_documento VARCHAR(2) NOT NULL,
			telefono VARCHAR(20) NOT NULL,
			direccion VARCHAR(150),
			correo_electronico VARCHAR(100)
		)`,

		`CREATE TABLE IF NOT EXISTS estudiantes_acudientes (
			id_estudiante INTEGER NOT NULL REFERENCES estudiantes(id_estudiante),
			id_acudiente INTEGER NOT NULL REFERENCES acudientes(id_acudiente),
			parentesco VARCHAR(50),
			es_principal BOOLEAN DEFAULT FALSE,
			PRIMARY KEY (id_estudiante, id_acudiente)
		)`,

		`CREATE TABLE IF NOT EXISTS roles (
			id_rol SERIAL PRIMARY KEY,
			nombre VARCHAR(50) NOT NULL UNIQUE,
			descripcion VARCHAR(255),
			es_rol_sistema BOOLEAN DEFAULT FALSE
		)`,

		`CREATE TABLE IF NOT EXISTS modulos (
			id_modulo SERIAL PRIMARY KEY,
			nombre VARCHAR(50) NOT NULL,
			descripcion VARCHAR(255),
			codigo VARCHAR(30) NOT NULL UNIQUE
		)`,

		`CREATE TABLE IF NOT EXISTS permisos (
			id_permiso SERIAL PRIMARY KEY,
			nombre VARCHAR(100) NOT NULL,
			descripcion VARCHAR(255),
			codigo VARCHAR(30) NOT NULL UNIQUE,
			id_modulo INTEGER NOT NULL REFERENCES modulos(id_modulo)
		)`,

		`CREATE TABLE IF NOT EXISTS roles_permisos (
			id_rol INTEGER NOT NULL REFERENCES roles(id_rol),
			id_permiso INTEGER NOT NULL REFERENCES permisos(id_permiso),
			PRIMARY KEY (id_rol, id_permiso)
		)`,

		`CREATE TABLE IF NOT EXISTS usuarios (
			id_usuario SERIAL PRIMARY KEY,
			nombre_usuario VARCHAR(50) NOT NULL UNIQUE,
			clave VARCHAR(255) NOT NULL,
			id_docente INTEGER REFERENCES docentes(id_docente),
			id_estudiante INTEGER REFERENCES estudiantes(id_estudiante),
			activo BOOLEAN DEFAULT TRUE,
			id_rol INTEGER NOT NULL REFERENCES roles(id_rol)
		)`,

		`CREATE TABLE IF NOT EXISTS actividades_recuperacion (
			id_recuperacion SERIAL PRIMARY KEY,
			id_calificacion INTEGER NOT NULL REFERENCES calificaciones(id_calificacion),
			fecha DATE NOT NULL,
			descripcion TEXT NOT NULL,
			nota_recuperacion DECIMAL(3,2),
			fecha_calificacion TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS dias_semana (
			id_dia SERIAL PRIMARY KEY,
			nombre VARCHAR(20) NOT NULL,
			codigo INTEGER NOT NULL UNIQUE
		)`,

		`CREATE TABLE IF NOT EXISTS horarios (
			id_horario SERIAL PRIMARY KEY,
			id_carga_academica INTEGER NOT NULL REFERENCES cargas_academicas(id_carga),
			id_dia INTEGER NOT NULL REFERENCES dias_semana(id_dia),
			hora_inicio TIME NOT NULL,
			hora_fin TIME NOT NULL,
			UNIQUE(id_carga_academica, id_dia, hora_inicio)
		)`,

		`CREATE TABLE IF NOT EXISTS archivos_digitales (
			id_archivo SERIAL PRIMARY KEY,
			tipo_archivo VARCHAR(50) NOT NULL,
			nombre_archivo VARCHAR(255) NOT NULL,
			ruta_almacenamiento VARCHAR(500) NOT NULL,
			mime_type VARCHAR(100) NOT NULL,
			fecha_carga TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			id_usuario_carga INTEGER REFERENCES usuarios(id_usuario),
			entidad_relacionada VARCHAR(50) NOT NULL,
			id_entidad INTEGER NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			token_hash VARCHAR(255) NOT NULL UNIQUE,
			jti VARCHAR(36) NOT NULL UNIQUE,
			id_usuario INTEGER NOT NULL REFERENCES usuarios(id_usuario) ON DELETE CASCADE,
			fecha_expiracion TIMESTAMP NOT NULL,
			fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			dispositivo VARCHAR(255),
			activo BOOLEAN DEFAULT TRUE
		)`,

		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_usuario ON refresh_tokens(id_usuario)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expiracion ON refresh_tokens(fecha_expiracion)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_jti ON refresh_tokens(jti)`,

		`CREATE TABLE IF NOT EXISTS revoked_access_tokens (
			jti VARCHAR(36) PRIMARY KEY,
			fecha_revocado TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			fecha_expiracion TIMESTAMP NOT NULL
		)`,

		`CREATE INDEX IF NOT EXISTS idx_revoked_expira ON revoked_access_tokens(fecha_expiracion)`,

		`INSERT INTO dias_semana (nombre, codigo) VALUES 
			('Lunes', 1), ('Martes', 2), ('Miércoles', 3), 
			('Jueves', 4), ('Viernes', 5), ('Sábado', 6), ('Domingo', 7)
		ON CONFLICT (codigo) DO NOTHING`,

		`INSERT INTO jornadas (nombre) VALUES 
			('Mañana'), ('Tarde'), ('Completa'), ('Sabatina')
		ON CONFLICT DO NOTHING`,

		`INSERT INTO roles (nombre, descripcion, es_rol_sistema) VALUES 
			('ADMIN', 'Administrador del sistema', TRUE),
			('DIRECTOR', 'Director de institución', TRUE),
			('DOCENTE', 'Docente', TRUE),
			('ESTUDIANTE', 'Estudiante', TRUE),
			('ACUDIENTE', 'Acudiente', TRUE)
		ON CONFLICT (nombre) DO NOTHING`,

		`INSERT INTO modulos (nombre, descripcion, codigo) VALUES 
			('Usuarios', 'Gestión de usuarios del sistema', 'USERS'),
			('Roles', 'Gestión de roles del sistema', 'ROLES'),
			('Permisos', 'Gestión de permisos del sistema', 'PERMS'),
			('Instituciones', 'Gestión de instituciones educativas', 'INST'),
			('Sedes', 'Gestión de sedes', 'SEDES'),
			('Grados', 'Gestión de grados', 'GRADOS'),
			('Grupos', 'Gestión de grupos', 'GRUPOS'),
			('Jornadas', 'Gestión de jornadas', 'JORN'),
			('Asignaturas', 'Gestión de asignaturas', 'ASIG'),
			('Estudiantes', 'Gestión de estudiantes', 'ESTU'),
			('Docentes', 'Gestión de docentes', 'DOCE'),
			('Matrículas', 'Gestión de matrículas', 'MATR'),
			('Períodos', 'Gestión de períodos académicos', 'PERI'),
			('Evaluaciones', 'Gestión de evaluaciones', 'EVAL'),
			('Calificaciones', 'Gestión de calificaciones', 'CALI'),
			('Cargas', 'Gestión de cargas académicas', 'CARG'),
			('Horarios', 'Gestión de horarios', 'HOR'),
			('Acudientes', 'Gestión de acudientes', 'ACUD'),
			('Reportes', 'Generación de reportes', 'REPT')
		ON CONFLICT (codigo) DO NOTHING`,

		`INSERT INTO permisos (nombre, descripcion, codigo, id_modulo) VALUES 
			-- Módulo Usuarios
			('Ver usuarios', 'Ver lista de usuarios', 'usuarios_ver', (SELECT id_modulo FROM modulos WHERE codigo = 'USERS')),
			('Crear usuario', 'Crear nuevos usuarios', 'usuarios_crear', (SELECT id_modulo FROM modulos WHERE codigo = 'USERS')),
			('Editar usuario', 'Editar usuarios existentes', 'usuarios_editar', (SELECT id_modulo FROM modulos WHERE codigo = 'USERS')),
			('Eliminar usuario', 'Eliminar usuarios', 'usuarios_eliminar', (SELECT id_modulo FROM modulos WHERE codigo = 'USERS')),
			('Cambiar contraseña', 'Cambiar contraseña de usuarios', 'usuarios_cambiar_clave', (SELECT id_modulo FROM modulos WHERE codigo = 'USERS')),
			-- Módulo Roles
			('Ver roles', 'Ver lista de roles', 'roles_ver', (SELECT id_modulo FROM modulos WHERE codigo = 'ROLES')),
			('Crear rol', 'Crear nuevos roles', 'roles_crear', (SELECT id_modulo FROM modulos WHERE codigo = 'ROLES')),
			('Editar rol', 'Editar roles existentes', 'roles_editar', (SELECT id_modulo FROM modulos WHERE codigo = 'ROLES')),
			('Eliminar rol', 'Eliminar roles', 'roles_eliminar', (SELECT id_modulo FROM modulos WHERE codigo = 'ROLES')),
			-- Módulo Permisos
			('Ver permisos', 'Ver lista de permisos', 'permisos_ver', (SELECT id_modulo FROM modulos WHERE codigo = 'PERMS')),
			-- Módulo Instituciones
			('Ver instituciones', 'Ver lista de instituciones', 'instituciones_ver', (SELECT id_modulo FROM modulos WHERE codigo = 'INST')),
			('Crear institución', 'Crear nuevas instituciones', 'instituciones_crear', (SELECT id_modulo FROM modulos WHERE codigo = 'INST')),
			('Editar institución', 'Editar instituciones existentes', 'instituciones_editar', (SELECT id_modulo FROM modulos WHERE codigo = 'INST')),
			('Eliminar institución', 'Eliminar instituciones', 'instituciones_eliminar', (SELECT id_modulo FROM modulos WHERE codigo = 'INST')),
			-- Módulo Sedes
			('Ver sedes', 'Ver lista de sedes', 'sedes_ver', (SELECT id_modulo FROM modulos WHERE codigo = 'SEDES')),
			('Crear sede', 'Crear nuevas sedes', 'sedes_crear', (SELECT id_modulo FROM modulos WHERE codigo = 'SEDES')),
			('Editar sede', 'Editar sedes existentes', 'sedes_editar', (SELECT id_modulo FROM modulos WHERE codigo = 'SEDES')),
			('Eliminar sede', 'Eliminar sedes', 'sedes_eliminar', (SELECT id_modulo FROM modulos WHERE codigo = 'SEDES')),
			-- Módulo Grados
			('Ver grados', 'Ver lista de grados', 'grados_ver', (SELECT id_modulo FROM modulos WHERE codigo = 'GRADOS')),
			('Crear grado', 'Crear nuevos grados', 'grados_crear', (SELECT id_modulo FROM modulos WHERE codigo = 'GRADOS')),
			('Editar grado', 'Editar grados existentes', 'grados_editar', (SELECT id_modulo FROM modulos WHERE codigo = 'GRADOS')),
			('Eliminar grado', 'Eliminar grados', 'grados_eliminar', (SELECT id_modulo FROM modulos WHERE codigo = 'GRADOS')),
			-- Módulo Estudiantes
			('Ver estudiantes', 'Ver lista de estudiantes', 'estudiantes_ver', (SELECT id_modulo FROM modulos WHERE codigo = 'ESTU')),
			('Crear estudiante', 'Crear nuevos estudiantes', 'estudiantes_crear', (SELECT id_modulo FROM modulos WHERE codigo = 'ESTU')),
			('Editar estudiante', 'Editar estudiantes existentes', 'estudiantes_editar', (SELECT id_modulo FROM modulos WHERE codigo = 'ESTU')),
			('Eliminar estudiante', 'Eliminar estudiantes', 'estudiantes_eliminar', (SELECT id_modulo FROM modulos WHERE codigo = 'ESTU')),
			-- Módulo Docentes
			('Ver docentes', 'Ver lista de docentes', 'docentes_ver', (SELECT id_modulo FROM modulos WHERE codigo = 'DOCE')),
			('Crear docente', 'Crear nuevos docentes', 'docentes_crear', (SELECT id_modulo FROM modulos WHERE codigo = 'DOCE')),
			('Editar docente', 'Editar docentes existentes', 'docentes_editar', (SELECT id_modulo FROM modulos WHERE codigo = 'DOCE')),
			('Eliminar docente', 'Eliminar docentes', 'docentes_eliminar', (SELECT id_modulo FROM modulos WHERE codigo = 'DOCE')),
			-- Módulo Matrículas
			('Ver matrículas', 'Ver lista de matrículas', 'matriculas_ver', (SELECT id_modulo FROM modulos WHERE codigo = 'MATR')),
			('Crear matrícula', 'Crear nuevas matrículas', 'matriculas_crear', (SELECT id_modulo FROM modulos WHERE codigo = 'MATR')),
			('Editar matrícula', 'Editar matrículas existentes', 'matriculas_editar', (SELECT id_modulo FROM modulos WHERE codigo = 'MATR')),
			('Eliminar matrícula', 'Eliminar matrículas', 'matriculas_eliminar', (SELECT id_modulo FROM modulos WHERE codigo = 'MATR'))
		ON CONFLICT (codigo) DO NOTHING`,

		`INSERT INTO roles_permisos (id_rol, id_permiso) 
			SELECT r.id_rol, p.id_permiso 
			FROM roles r, permisos p 
			WHERE r.nombre = 'ADMIN'
		ON CONFLICT DO NOTHING`,

		`INSERT INTO usuarios (nombre_usuario, clave, activo, id_rol) VALUES 
			('admin', '$argon2id$v=19$m=65536,t=1,p=8$Q6wV+rU8BSyEFkAVmUVTHA$jQeYp/eTbOqRYO8FghFOt5b9Q5K8PLj3BXPHti6vPnc', true, 1)
		ON CONFLICT (nombre_usuario) DO NOTHING`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w\nQuery: %s", err, migration)
		}
	}

	return nil
}
