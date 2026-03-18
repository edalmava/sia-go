document.addEventListener('DOMContentLoaded', function() {
    // Check authentication
    const token = localStorage.getItem('authToken');
    if (!token) {
        window.location.href = 'login.html';
        return;
    }

    // Verify token and load permissions
    try {
        const decoded = jwt_decode(token);
        const permisos = decoded.permisos || [];
        localStorage.setItem('userPermissions', JSON.stringify(permisos));
        localStorage.setItem('userRole', decoded.rol || 'USER');
    } catch (e) {
        console.error('Error decoding token:', e);
    }

    // Load user info
    const username = localStorage.getItem('username') || 'admin';
    const userRole = localStorage.getItem('userRole') || 'USER';
    const permisos = JSON.parse(localStorage.getItem('userPermissions') || '[]');
    
    document.getElementById('userName').textContent = username.charAt(0).toUpperCase() + username.slice(1);
    document.getElementById('userInitial').textContent = username.charAt(0).toUpperCase();
    
    // Set role text based on user role
    const roleNames = {
        'ADMIN': 'Administrador',
        'DIRECTOR': 'Director',
        'DOCENTE': 'Docente',
        'ESTUDIANTE': 'Estudiante',
        'ACUDIENTE': 'Acudiente',
        'USER': 'Usuario'
    };
    const roleElement = document.querySelector('.user-role');
    if (roleElement) {
        roleElement.textContent = roleNames[userRole] || 'Usuario';
    }

    // Show menu items based on permissions (new data-permiso approach)
    document.querySelectorAll('.nav-item[data-permiso]').forEach(item => {
        const permiso = item.dataset.permiso;
        if (permisos.includes(permiso)) {
            item.style.display = 'flex';
        }
    });

    // Logout handler
    document.getElementById('logoutBtn').addEventListener('click', function(e) {
        e.preventDefault();
        localStorage.removeItem('authToken');
        localStorage.removeItem('username');
        localStorage.removeItem('userRole');
        localStorage.removeItem('userPermissions');
        window.location.href = 'login.html';
    });

    // Active nav item handler
    const navItems = document.querySelectorAll('.nav-item');
    navItems.forEach(item => {
        item.addEventListener('click', function(e) {
            const href = this.getAttribute('href');
            // Allow navigation to other pages
            if (href && href !== '#') {
                // Just add active class, let navigation happen naturally
                navItems.forEach(nav => nav.classList.remove('active'));
                this.classList.add('active');
                return; // Don't prevent default
            }
            // For # links, prevent default
            e.preventDefault();
            navItems.forEach(nav => nav.classList.remove('active'));
            this.classList.add('active');
        });
    });

    // Fetch dashboard data (mock for now)
    loadDashboardData();
});

async function loadDashboardData() {
    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch('/api/v1/instituciones', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (response.ok) {
            const data = await response.json();
            console.log('Dashboard data loaded:', data);
        }
    } catch (error) {
        console.log('Using mock data for dashboard');
    }
}