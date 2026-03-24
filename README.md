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
- **Interfaz Moderna:** Dashboard responsivo construido con Vanilla JavaScript y CSS3, optimizado con **Vite**.
- **API RESTful:** Documentación completa siguiendo la especificación OpenAPI 3.0.

## 🛠️ Tecnologías Utilizadas

- **Backend:** [Go](https://go.dev/) (v1.24.1) con el framework [Echo](https://echo.labstack.com/).
- **Base de Datos:** PostgreSQL (usando `lib/pq`).
- **Frontend:** HTML5, CSS3, JavaScript Vanilla y **Vite** para optimización.
- **Gestión de Entorno:** [godotenv](https://github.com/joho/godotenv) para carga automática de archivos `.env`.
- **Autenticación:** JWT (JSON Web Tokens).

## 📁 Estructura del Proyecto

- `cmd/server/`: Punto de entrada de la aplicación y configuración del servidor.
- `internal/handlers/`: Controladores de solicitudes HTTP.
- `internal/repository/`: Capa de acceso a datos (consultas SQL).
- `internal/models/`: Definiciones de estructuras de datos y entidades de la DB.
- `src-web/`: **Código fuente del frontend** (archivos de desarrollo).
- `web/`: **Distribución del frontend** (archivos optimizados generados por Vite).
- `api/`: Definición de la API en formato OpenAPI.

## ⚙️ Configuración y Ejecución

### Requisitos Previos

- Go 1.24 o superior.
- Node.js (v18+) y npm para el desarrollo del frontend.
- Base de datos PostgreSQL.
- Archivo `.env` (ver `.env.example`).

### Gestión del Frontend (Vite)

El frontend utiliza Vite para agrupar (bundle) y minimizar los activos, reduciendo el número de peticiones HTTP y mejorando el rendimiento.

- **Instalar dependencias:** `npm install`
- **Modo desarrollo:** `npm run dev` (Inicia el servidor de desarrollo de Vite).
- **Construir para producción:** `npm run build` (Genera los archivos optimizados en la carpeta `web/`).

*Nota: El servidor Go sirve los archivos desde la carpeta `web/`. Siempre ejecute `npm run build` después de modificar archivos en `src-web/` para ver los cambios reflejados en el servidor de backend.*

### Pasos para Ejecutar el Sistema

1. **Configurar el Entorno:**
   ```bash
   cp .env.example .env
   # Edite .env con sus credenciales de base de datos y un JWT_SECRET seguro.
   ```

2. **Construir el Frontend:**
   ```bash
   npm install
   npm run build
   ```

3. **Iniciar el Servidor Go:**
   ```bash
   go run cmd/server/main.go
   ```
   El servidor se iniciará por defecto en el puerto `8080`.

3. **Acceder a la Aplicación:**
   - **Interfaz Web:** [http://localhost:8080/](http://localhost:8080/)
   - **Salud de la API:** `GET /health`

## 📄 Licencia

Este proyecto es para uso académico e institucional. Todos los derechos reservados.
