import { AUTH_API } from './config.js';
import { api } from './api.js';
import { isLoggedIn } from './auth.js';
import { showToast } from './ui.js';

document.addEventListener('DOMContentLoaded', () => {
    if (isLoggedIn()) {
        window.location.href = 'pages/dashboard.html';
        return;
    }

    initLoginForm();
    initPasswordToggle();
});

function initLoginForm() {
    const loginForm = document.getElementById('loginForm');
    const loginBtn = document.getElementById('loginBtn');
    const errorMessage = document.getElementById('errorMessage');

    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        errorMessage.classList.remove('show');
        
        const username = document.getElementById('username').value.trim();
        const password = document.getElementById('password').value;

        if (!username || !password) {
            showError('Por favor ingresa usuario y contraseña');
            return;
        }

        loginBtn.classList.add('loading');
        loginBtn.disabled = true;

        try {
            const data = await api.postAuth(AUTH_API + '/login', { username, password });

            if (data.access_token) {
                saveUserData(data);
                window.location.href = 'pages/dashboard.html';
            } else {
                showError(data.message || 'Credenciales inválidas');
            }
        } catch (error) {
            showError(error.data?.message || 'Credenciales inválidas');
        } finally {
            loginBtn.classList.remove('loading');
            loginBtn.disabled = false;
        }
    });
}

function saveUserData(data) {
    localStorage.setItem('authToken', data.access_token);
    localStorage.setItem('username', data.nombre_usuario);
    localStorage.setItem('userRole', data.role || 'USER');
    localStorage.setItem('idRol', data.id_rol);
    localStorage.setItem('userPermissions', JSON.stringify(data.permisos || []));
}

function showError(message) {
    const errorMessage = document.getElementById('errorMessage');
    errorMessage.querySelector('span').textContent = message;
    errorMessage.classList.add('show');
    
    setTimeout(() => {
        errorMessage.classList.remove('show');
    }, 5000);
}

function initPasswordToggle() {
    const toggleBtn = document.getElementById('togglePassword');
    const passwordInput = document.getElementById('password');

    if (!toggleBtn || !passwordInput) return;

    toggleBtn.addEventListener('click', () => {
        const type = passwordInput.getAttribute('type') === 'password' ? 'text' : 'password';
        passwordInput.setAttribute('type', type);
        
        const eyeIcon = toggleBtn.querySelector('.eye-icon');
        if (type === 'text') {
            eyeIcon.innerHTML = `
                <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                <line x1="1" y1="1" x2="23" y2="23"/>
            `;
        } else {
            eyeIcon.innerHTML = `
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
            `;
        }
    });
}
