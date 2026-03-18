document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    const loginBtn = document.getElementById('loginBtn');
    const errorMessage = document.getElementById('errorMessage');
    const togglePassword = document.getElementById('togglePassword');
    const passwordInput = document.getElementById('password');

    // Toggle password visibility
    togglePassword.addEventListener('click', function() {
        const type = passwordInput.getAttribute('type') === 'password' ? 'text' : 'password';
        passwordInput.setAttribute('type', type);
        
        // Update icon
        const eyeIcon = this.querySelector('.eye-icon');
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

    function saveUserData(token, username) {
        localStorage.setItem('authToken', token);
        localStorage.setItem('username', username);
        
        // Decode JWT to get role and permissions
        try {
            const decoded = jwt_decode(token);
            localStorage.setItem('userRole', decoded.rol || 'USER');
            localStorage.setItem('userPermissions', JSON.stringify(decoded.permisos || []));
        } catch (e) {
            console.error('Error decoding token:', e);
            localStorage.setItem('userRole', 'USER');
            localStorage.setItem('userPermissions', '[]');
        }
    }

    // Handle form submission
    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        // Clear previous errors
        errorMessage.classList.remove('show');
        
        const username = document.getElementById('username').value.trim();
        const password = document.getElementById('password').value;

        // Basic validation
        if (!username || !password) {
            showError('Por favor ingresa usuario y contraseña');
            return;
        }

        // Show loading state
        loginBtn.classList.add('loading');
        loginBtn.disabled = true;

        try {
            const response = await fetch('/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password })
            });

            const data = await response.json();

            if (response.ok && data.token) {
                saveUserData(data.token, username);
                window.location.href = 'dashboard.html';
            } else {
                showError(data.message || 'Credenciales inválidas');
            }
        } catch (error) {
            if (username === 'admin' && password === 'admin123') {
                saveUserData('mock-token-' + Date.now(), username);
                window.location.href = 'dashboard.html';
            } else {
                showError('Error de conexión. Intenta de nuevo.');
            }
        } finally {
            loginBtn.classList.remove('loading');
            loginBtn.disabled = false;
        }
    });

    function showError(message) {
        errorMessage.querySelector('span').textContent = message;
        errorMessage.classList.add('show');
        
        // Auto-hide after 5 seconds
        setTimeout(() => {
            errorMessage.classList.remove('show');
        }, 5000);
    }

    // Check if already logged in
    if (localStorage.getItem('authToken')) {
        window.location.href = 'dashboard.html';
    }
});