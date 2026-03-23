# SIA - Sistema Integrado Académico

SIA es un sistema integral de gestión escolar diseñado para instituciones educativas en Colombia. Gestiona desde datos institucionales básicos hasta estructuras académicas complejas como grados, grupos, estudiantes, docentes, matrículas y calificaciones.

## 🚀 Características Principales

- **Gestión Institucional:** Administración de sedes, jornadas y configuración del establecimiento educativo.
- **Estructura Académica:** Control de grados, grupos y asignaturas.
- **Gestión de Usuarios:** Sistema de roles y permisos granulares.
- **Seguridad Avanzada:**
  - Autenticación mediante JWT con rotación de tokens.
  - Almacenamiento seguro de tokens de refresco en cookies `HttpOnly`.
  - Protección de contraseñas con Argon2id.
  - Implementación de políticas de seguridad (CSP, HSTS, CORS restringido).
- **Interfaz Intuitiva:** Dashboard moderno construido con tecnologías web estándar (HTML5, CSS3, JavaScript Vanilla).
- **API RESTful:** Documentación completa siguiendo la especificación OpenAPI 3.0.

## 🛠️ Tecnologías Utilizadas

- **Backend:** [Go](https://go.dev/) (v1.24.1) con el framework [Echo](https://echo.labstack.com/).
- **Base de Datos:** PostgreSQL (usando `lib/pq`).
- **Frontend:** HTML5, CSS3 y JavaScript nativo (sin frameworks pesados).
- **Gestión de Entorno:** [godotenv](https://github.com/joho/godotenv) para carga automática de archivos `.env`.
- **Autenticación:** JWT (JSON Web Tokens).

## 📁 Estructura del Proyecto

- `cmd/server/`: Punto de entrada de la aplicación y configuración del servidor.
- `internal/handlers/`: Controladores de solicitudes HTTP.
- `internal/repository/`: Capa de acceso a datos (consultas SQL).
- `internal/models/`: Definiciones de estructuras de datos y entidades de la DB.
- `internal/database/`: Conexión a la base de datos y migraciones automáticas.
- `internal/middleware/`: Middleware de autenticación, seguridad y registro.
- `web/`: Archivos estáticos del frontend (HTML, CSS, JS).
- `api/`: Definición de la API en formato OpenAPI.

## ⚙️ Configuración y Ejecución

### Requisitos Previos

- Go 1.24 o superior.
- Base de datos PostgreSQL.
- Archivo `.env` (ver `.env.example`).

### Pasos para Ejecutar

1. **Configurar el Entorno:**
   ```bash
   cp .env.example .env
   # Edite .env con sus credenciales de base de datos y un JWT_SECRET seguro.
   ```

2. **Iniciar el Servidor:**
   ```bash
   go run cmd/server/main.go
   ```
   El servidor se iniciará por defecto en el puerto `8080`.

3. **Acceder a la Aplicación:**
   - **Interfaz Web:** [http://localhost:8080/](http://localhost:8080/)
   - **Salud de la API:** `GET /health`

## 📄 Licencia

Este proyecto es para uso académico e institucional. Todos los derechos reservados.
