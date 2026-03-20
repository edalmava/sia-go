// Configuration constants
export const API_URL = '/api/v1';
export const AUTH_API = '/auth';

export const ROLE_NAMES = {
    'ADMIN': 'Administrador',
    'DIRECTOR': 'Director',
    'DOCENTE': 'Docente',
    'ESTUDIANTE': 'Estudiante',
    'ACUDIENTE': 'Acudiente',
    'USER': 'Usuario'
};

export const ROLE_COLORS = {
    'ADMIN': 'admin',
    'DIRECTOR': 'director',
    'DOCENTE': 'docente',
    'ESTUDIANTE': 'estudiante',
    'ACUDIENTE': 'acudiente',
    'USER': 'user'
};

export const PAGINATION = {
    DEFAULT_LIMIT: 20,
    MAX_LIMIT: 100
};

export const API_ENDPOINTS = {
    USUARIOS: '/usuarios',
    ROLES: '/roles',
    PERMISOS: '/permisos',
    MODULOS: '/modulos',
    INSTITUCIONES: '/instituciones',
    SEDES: '/sedes',
    GRADOS: '/grados',
    GRUPOS: '/grupos',
    JORNADAS: '/jornadas',
    ASIGNATURAS: '/asignaturas',
    ESTUDIANTES: '/estudiantes',
    DOCENTES: '/docentes'
};
