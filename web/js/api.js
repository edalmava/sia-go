import { API_URL, AUTH_API } from './config.js';

const STORAGE_KEYS = {
    TOKEN: 'authToken'
};

async function getHeaders() {
    return {
        'Content-Type': 'application/json'
    };
}

async function refreshAccessToken() {
    try {
        const response = await fetch(`${AUTH_API}/refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include'
        });

        if (!response.ok) {
            return false;
        }

        const data = await response.json();
        
        // El access_token ahora viene en una cookie HttpOnly
        // Actualizar datos del usuario que podrían haber cambiado
        localStorage.setItem('username', data.nombre_usuario);
        localStorage.setItem('userRole', data.role || 'USER');
        localStorage.setItem('idRol', data.id_rol);
        localStorage.setItem('userPermissions', JSON.stringify(data.permisos || []));
        // Guardamos un flag para saber que estamos logueados ya que no podemos leer la cookie HttpOnly
        localStorage.setItem('isLoggedIn', 'true');
        return true;
    } catch (e) {
        console.error('Error refreshing token:', e);
    }
    return false;
}

function clearAuthAndRedirect() {
    localStorage.removeItem('isLoggedIn');
    localStorage.removeItem('username');
    localStorage.removeItem('userRole');
    localStorage.removeItem('idRol');
    localStorage.removeItem('userPermissions');
    
    window.location.href = '/web/login.html';
}

async function handleResponse(response, retry = true) {
    const text = await response.text();
    
    if (!response.ok) {
        if (response.status === 401 && retry) {
            const refreshed = await refreshAccessToken();
            if (refreshed) {
                window.location.reload(); 
                return;
            }
            clearAuthAndRedirect();
            return;
        }

        let message = 'Error en la solicitud';
        let data = null;
        
        try {
            if (text) {
                data = JSON.parse(text);
                message = data.message || message;
            }
        } catch (e) {
            console.error('Error parsing error response:', e);
        }

        const error = new Error(message);
        error.status = response.status;
        error.data = data;
        throw error;
    }
    
    if (!text) return null;
    return JSON.parse(text);
}

export const api = {
    async get(endpoint, params = {}) {
        const queryString = new URLSearchParams(params).toString();
        const url = queryString ? `${API_URL}${endpoint}?${queryString}` : `${API_URL}${endpoint}`;
        
        const response = await fetch(url, {
            method: 'GET',
            headers: await getHeaders(),
            credentials: 'include'
        });
        
        return handleResponse(response);
    },

    async post(endpoint, data) {
        const response = await fetch(`${API_URL}${endpoint}`, {
            method: 'POST',
            headers: await getHeaders(),
            body: JSON.stringify(data),
            credentials: 'include'
        });
        
        return handleResponse(response);
    },

    async put(endpoint, data) {
        const response = await fetch(`${API_URL}${endpoint}`, {
            method: 'PUT',
            headers: await getHeaders(),
            body: JSON.stringify(data),
            credentials: 'include'
        });
        
        return handleResponse(response);
    },

    async delete(endpoint) {
        const response = await fetch(`${API_URL}${endpoint}`, {
            method: 'DELETE',
            headers: await getHeaders(),
            credentials: 'include'
        });
        
        return handleResponse(response);
    },

    async postAuth(endpoint, data) {
        const response = await fetch(`${endpoint}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
            credentials: 'include'
        });
        
        return handleResponse(response);
    }
};

// Specific API modules
export const usuarioApi = {
    getAll: (offset = 0, limit = 20, search = '') => 
        api.get('/usuarios', { offset, limit, search }),
    getById: (id) => api.get(`/usuarios/${id}`),
    create: (data) => api.post('/usuarios', data),
    update: (id, data) => api.put(`/usuarios/${id}`, data),
    delete: (id) => api.delete(`/usuarios/${id}`),
    changePassword: (id, newPassword) => api.post(`/usuarios/${id}/change-password`, { password: newPassword })
};

export const rolApi = {
    getAll: () => api.get('/roles'),
    getById: (id) => api.get(`/roles/${id}`),
    create: (data) => api.post('/roles', data),
    update: (id, data) => api.put(`/roles/${id}`, data),
    delete: (id) => api.delete(`/roles/${id}`)
};

export const permisoApi = {
    getAll: () => api.get('/permisos'),
    getById: (id) => api.get(`/permisos/${id}`)
};

export const moduloApi = {
    getAll: () => api.get('/modulos'),
    getById: (id) => api.get(`/modulos/${id}`)
};

export const institucionApi = {
    getAll: (offset = 0, limit = 20) => api.get('/instituciones', { offset, limit }),
    getById: (id) => api.get(`/instituciones/${id}`)
};
