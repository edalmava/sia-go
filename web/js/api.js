import { API_URL } from './config.js';

async function getHeaders() {
    const token = localStorage.getItem('authToken');
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    };
}

async function handleResponse(response) {
    const text = await response.text();
    
    if (!response.ok) {
        if (response.status === 401) {
            localStorage.removeItem('authToken');
            localStorage.removeItem('userData');
            window.location.href = '/web/login.html';
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
            headers: await getHeaders()
        });
        
        return handleResponse(response);
    },

    async post(endpoint, data) {
        const response = await fetch(`${API_URL}${endpoint}`, {
            method: 'POST',
            headers: await getHeaders(),
            body: JSON.stringify(data)
        });
        
        return handleResponse(response);
    },

    async put(endpoint, data) {
        const response = await fetch(`${API_URL}${endpoint}`, {
            method: 'PUT',
            headers: await getHeaders(),
            body: JSON.stringify(data)
        });
        
        return handleResponse(response);
    },

    async delete(endpoint) {
        const response = await fetch(`${API_URL}${endpoint}`, {
            method: 'DELETE',
            headers: await getHeaders()
        });
        
        return handleResponse(response);
    },

    async postAuth(endpoint, data) {
        const response = await fetch(`${endpoint}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
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
