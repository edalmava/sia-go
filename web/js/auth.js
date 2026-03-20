import { ROLE_NAMES } from './config.js';

const STORAGE_KEYS = {
    TOKEN: 'authToken',
    USERNAME: 'username',
    ROLE: 'userRole',
    PERMISSIONS: 'userPermissions'
};

function getLoginPath() {
    const currentPath = window.location.pathname;
    if (currentPath.includes('/pages/')) {
        return '../login.html';
    }
    return 'login.html';
}

export function checkAuth(redirectTo = null) {
    const token = localStorage.getItem(STORAGE_KEYS.TOKEN);
    if (!token) {
        if (redirectTo) {
            window.location.href = redirectTo;
        } else {
            window.location.href = getLoginPath();
        }
        return false;
    }

    try {
        const decoded = jwt_decode(token);
        localStorage.setItem(STORAGE_KEYS.PERMISSIONS, JSON.stringify(decoded.permisos || []));
        localStorage.setItem(STORAGE_KEYS.ROLE, decoded.rol || 'USER');
    } catch (e) {
        console.error('Error decoding token:', e);
    }

    return true;
}

export function requireAuth() {
    if (!checkAuth()) {
        return false;
    }
    return true;
}

export function getUserData() {
    return {
        username: localStorage.getItem(STORAGE_KEYS.USERNAME) || 'admin',
        role: localStorage.getItem(STORAGE_KEYS.ROLE) || 'USER',
        permissions: JSON.parse(localStorage.getItem(STORAGE_KEYS.PERMISSIONS) || '[]')
    };
}

export function hasPermission(permiso) {
    const { permissions } = getUserData();
    return permissions.includes(permiso);
}

export function hasAnyPermission(permisos) {
    const { permissions } = getUserData();
    return permisos.some(p => permissions.includes(p));
}

export function hasPermissionPrefix(prefix) {
    const { permissions } = getUserData();
    return permissions.some(p => p.startsWith(prefix));
}

export function hasAnyPermissionPrefix(prefixes) {
    const { permissions } = getUserData();
    return prefixes.some(prefix => permissions.some(p => p.startsWith(prefix)));
}

export function logout(redirectTo = null) {
    localStorage.removeItem(STORAGE_KEYS.TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USERNAME);
    localStorage.removeItem(STORAGE_KEYS.ROLE);
    localStorage.removeItem(STORAGE_KEYS.PERMISSIONS);
    
    if (redirectTo) {
        window.location.href = redirectTo;
    } else {
        const currentPath = window.location.pathname;
        if (currentPath.includes('/pages/')) {
            window.location.href = '../login.html';
        } else {
            window.location.href = 'login.html';
        }
    }
}

export function initUserInfo() {
    const { username, role } = getUserData();
    
    const userNameEl = document.getElementById('userName');
    const userInitialEl = document.getElementById('userInitial');
    const userRoleEl = document.querySelector('.user-role');

    if (userNameEl) {
        userNameEl.textContent = username.charAt(0).toUpperCase() + username.slice(1);
    }
    
    if (userInitialEl) {
        userInitialEl.textContent = username.charAt(0).toUpperCase();
    }
    
    if (userRoleEl) {
        userRoleEl.textContent = ROLE_NAMES[role] || 'Usuario';
    }
}

export function initLogout() {
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', (e) => {
            e.preventDefault();
            logout();
        });
    }
}

export function initSidebar() {
    const { permissions } = getUserData();
    
    document.querySelectorAll('.nav-item[data-permiso]').forEach(item => {
        const permiso = item.dataset.permiso;
        if (permissions.includes(permiso)) {
            item.style.display = 'flex';
        }
    });
}

export function isLoggedIn() {
    return !!localStorage.getItem(STORAGE_KEYS.TOKEN);
}

export { STORAGE_KEYS, ROLE_NAMES };
