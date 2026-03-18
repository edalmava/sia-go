document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('authToken');
    if (!token) {
        window.location.href = 'login.html';
        return;
    }

    try {
        const decoded = jwt_decode(token);
        const permisos = decoded.permisos || [];
        localStorage.setItem('userPermissions', JSON.stringify(permisos));
        localStorage.setItem('userRole', decoded.rol || 'USER');
    } catch (e) {
        console.error('Error decoding token:', e);
    }

    const permisos = JSON.parse(localStorage.getItem('userPermissions') || '[]');
    
    // Show menu items based on permissions
    document.querySelectorAll('.nav-item[data-permiso]').forEach(item => {
        const permiso = item.dataset.permiso;
        if (permisos.includes(permiso)) {
            item.style.display = 'flex';
        }
    });

    // Check if user has permission to access settings (usuarios, roles, or permisos permissions)
    const canAccessSettings = permisos.some(p => 
        p.startsWith('usuarios_') || 
        p.startsWith('roles_') || 
        p.startsWith('permisos_')
    );
    
    if (!canAccessSettings) {
        showToast('No tienes permisos para acceder a esta página', 'error');
        setTimeout(() => {
            window.location.href = 'dashboard.html';
        }, 1500);
        return;
    }

    const username = localStorage.getItem('username') || 'admin';
    const userRole = localStorage.getItem('userRole') || 'USER';
    
    document.getElementById('userName').textContent = username.charAt(0).toUpperCase() + username.slice(1);
    document.getElementById('userInitial').textContent = username.charAt(0).toUpperCase();
    
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

    const tabButtons = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-content');
    
    // Show/hide tabs based on permissions
    const hasRolesPerms = permisos.some(p => p.startsWith('roles_'));
    const hasPermisosPerms = permisos.includes('permisos_ver');
    
    tabButtons.forEach(btn => {
        const tabName = btn.dataset.tab;
        
        // Hide tab if user doesn't have permission
        if (tabName === 'roles' && !hasRolesPerms) {
            btn.style.display = 'none';
            return;
        }
        if (tabName === 'permisos' && !hasPermisosPerms) {
            btn.style.display = 'none';
            return;
        }
        
        btn.addEventListener('click', () => {
            tabButtons.forEach(b => b.classList.remove('active'));
            tabContents.forEach(c => c.classList.remove('active'));
            btn.classList.add('active');
            document.getElementById(`${btn.dataset.tab}-tab`).classList.add('active');
        });
    });

    // Hide tab contents that user doesn't have permission to view
    if (!hasRolesPerms) {
        document.getElementById('roles-tab').style.display = 'none';
    }
    if (!hasPermisosPerms) {
        document.getElementById('permisos-tab').style.display = 'none';
    }

    // If first visible tab is not "roles", switch to it
    const visibleTabs = Array.from(tabButtons).filter(btn => btn.style.display !== 'none');
    if (visibleTabs.length > 0 && !document.querySelector('.tab-btn.active')?.dataset.tab) {
        visibleTabs[0].classList.add('active');
        const firstVisibleTabName = visibleTabs[0].dataset.tab;
        document.getElementById(`${firstVisibleTabName}-tab`).classList.add('active');
    }

    document.getElementById('logoutBtn').addEventListener('click', function(e) {
        e.preventDefault();
        localStorage.removeItem('authToken');
        localStorage.removeItem('username');
        localStorage.removeItem('userRole');
        localStorage.removeItem('userPermissions');
        window.location.href = 'login.html';
    });

    // Show/hide "Nuevo Rol" button based on permissions
    const canCreateRoles = permisos.includes('roles_crear');
    const addRoleBtn = document.getElementById('addRoleBtn');
    if (addRoleBtn) {
        addRoleBtn.style.display = canCreateRoles ? 'inline-flex' : 'none';
    }

    // Show/hide "Nuevo Permiso" button based on permissions (only admin)
    const addPermisoBtn = document.getElementById('addPermisoBtn');
    if (addPermisoBtn) {
        addPermisoBtn.style.display = 'none'; // Hide by default
    }

    loadRoles();
    loadPermisos();
    loadModulos();

    document.getElementById('addRoleBtn').addEventListener('click', () => openRoleModal());
    document.getElementById('closeRoleModal').addEventListener('click', closeRoleModal);
    document.getElementById('cancelRoleBtn').addEventListener('click', closeRoleModal);
    document.getElementById('roleModal').addEventListener('click', function(e) {
        if (e.target === this) closeRoleModal();
    });
    document.getElementById('roleForm').addEventListener('submit', handleRoleSubmit);

    document.getElementById('permisoSearch').addEventListener('input', filterPermisos);
    document.getElementById('moduloFilter').addEventListener('change', filterPermisos);
});

const API_URL = '/api/v1';

let allPermisos = [];
let allModulos = [];

async function loadRoles() {
    const tbody = document.getElementById('rolesTableBody');
    tbody.innerHTML = '<tr><td colspan="6" class="loading-cell">Cargando roles...</td></tr>';

    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`${API_URL}/roles`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) throw new Error('Error al cargar roles');

        const data = await response.json();
        renderRoles(data.data || []);
    } catch (error) {
        console.error('Error:', error);
        tbody.innerHTML = '<tr><td colspan="6" class="empty-cell">Error al cargar roles</td></tr>';
    }
}

function renderRoles(roles) {
    const tbody = document.getElementById('rolesTableBody');
    const permisos = JSON.parse(localStorage.getItem('userPermissions') || '[]');
    const canEditRoles = permisos.includes('roles_editar');
    const canDeleteRoles = permisos.includes('roles_eliminar');
    
    if (roles.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" class="empty-cell">No hay roles registrados</td></tr>';
        return;
    }

    tbody.innerHTML = roles.map(rol => {
        let actions = '';
        
        // Edit button - show if user has permission and role is not a system role (or allow editing system role name)
        if (canEditRoles) {
            actions += `
                <button class="btn-action edit" onclick="openRoleModal(${rol.id_rol}, '${escapeHtml(rol.nombre)}', '${escapeHtml(rol.descripcion || '')}', ${rol.es_rol_sistema})" title="Editar rol">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                </button>`;
        }
        
        // Delete button - show if user has permission and role is not a system role
        if (canDeleteRoles && !rol.es_rol_sistema) {
            actions += `
                <button class="btn-action delete" onclick="deleteRole(${rol.id_rol}, '${escapeHtml(rol.nombre)}')" title="Eliminar rol">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    </svg>
                </button>`;
        }
        
        return `
            <tr>
                <td>${rol.id_rol}</td>
                <td><strong>${escapeHtml(rol.nombre)}</strong></td>
                <td>${escapeHtml(rol.descripcion || '-')}</td>
                <td><span class="system-badge ${rol.es_rol_sistema ? 'yes' : 'no'}">${rol.es_rol_sistema ? 'Sí' : 'No'}</span></td>
                <td><span class="permissions-count">${rol.permisos_count || 0} permisos</span></td>
                <td>
                    <div class="action-buttons">
                        ${actions}
                    </div>
                </td>
            </tr>
        `;
    }).join('');
}

async function loadPermisos() {
    const tbody = document.getElementById('permisosTableBody');
    tbody.innerHTML = '<tr><td colspan="5" class="loading-cell">Cargando permisos...</td></tr>';

    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`${API_URL}/permisos`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) throw new Error('Error al cargar permisos');

        const data = await response.json();
        allPermisos = data.data || [];
        
        await loadModulosForFilter();
        renderPermisos(allPermisos);
    } catch (error) {
        console.error('Error:', error);
        tbody.innerHTML = '<tr><td colspan="5" class="empty-cell">Error al cargar permisos</td></tr>';
    }
}

async function loadModulosForFilter() {
    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`${API_URL}/modulos`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (response.ok) {
            const data = await response.json();
            allModulos = data.data || [];
            
            const select = document.getElementById('moduloFilter');
            allModulos.forEach(modulo => {
                const option = document.createElement('option');
                option.value = modulo.id_modulo;
                option.textContent = modulo.nombre;
                select.appendChild(option);
            });
        }
    } catch (error) {
        console.error('Error loading modulos:', error);
    }
}

function renderPermisos(permisos) {
    const tbody = document.getElementById('permisosTableBody');
    
    if (permisos.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" class="empty-cell">No hay permisos registrados</td></tr>';
        return;
    }

    tbody.innerHTML = permisos.map(permiso => `
        <tr>
            <td>${permiso.id_permiso}</td>
            <td><strong>${escapeHtml(permiso.nombre)}</strong></td>
            <td><code>${escapeHtml(permiso.codigo)}</code></td>
            <td>${escapeHtml(permiso.modulo_nombre || '-')}</td>
            <td>${escapeHtml(permiso.descripcion || '-')}</td>
        </tr>
    `).join('');
}

function filterPermisos() {
    const search = document.getElementById('permisoSearch').value.toLowerCase();
    const moduloId = document.getElementById('moduloFilter').value;

    let filtered = allPermisos;
    
    if (search) {
        filtered = filtered.filter(p => 
            p.nombre.toLowerCase().includes(search) ||
            p.codigo.toLowerCase().includes(search) ||
            (p.descripcion && p.descripcion.toLowerCase().includes(search))
        );
    }
    
    if (moduloId) {
        filtered = filtered.filter(p => p.id_modulo == moduloId);
    }

    renderPermisos(filtered);
}

async function loadModulos() {
    const grid = document.getElementById('modulosGrid');
    grid.innerHTML = '<div class="loading-card">Cargando módulos...</div>';

    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`${API_URL}/modulos`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) throw new Error('Error al cargar módulos');

        const data = await response.json();
        renderModulos(data.data || []);
    } catch (error) {
        console.error('Error:', error);
        grid.innerHTML = '<div class="empty-state">Error al cargar módulos</div>';
    }
}

function renderModulos(modulos) {
    const grid = document.getElementById('modulosGrid');
    
    if (modulos.length === 0) {
        grid.innerHTML = '<div class="empty-state">No hay módulos registrados</div>';
        return;
    }

    grid.innerHTML = modulos.map(modulo => `
        <div class="module-card">
            <div class="module-header">
                <div class="module-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                        <line x1="3" y1="9" x2="21" y2="9"/>
                        <line x1="9" y1="21" x2="9" y2="9"/>
                    </svg>
                </div>
                <div>
                    <div class="module-name">${escapeHtml(modulo.nombre)}</div>
                    <div class="module-code">${escapeHtml(modulo.codigo)}</div>
                </div>
            </div>
            <p class="module-desc">${escapeHtml(modulo.descripcion || 'Sin descripción')}</p>
            <div class="module-stats">
                <span class="module-stat">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                        <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                        <polyline points="22 4 12 14.01 9 11.01"/>
                    </svg>
                    ${modulo.permisos_count || 0} permisos
                </span>
            </div>
        </div>
    `).join('');
}

let editingRoleId = null;
let allPermissionsData = [];

async function openRoleModal(roleId = null, nombre = '', descripcion = '', esSistema = false) {
    const modal = document.getElementById('roleModal');
    const form = document.getElementById('roleForm');
    const title = document.getElementById('roleModalTitle');

    editingRoleId = roleId;
    form.reset();

    if (roleId) {
        title.textContent = 'Editar Rol';
        document.getElementById('roleId').value = roleId;
        document.getElementById('roleName').value = nombre;
        document.getElementById('roleDesc').value = descripcion;
        document.getElementById('roleSistema').checked = esSistema;
        document.getElementById('roleName').disabled = esSistema;
        
        await loadPermissionsForRole(roleId);
    } else {
        title.textContent = 'Nuevo Rol';
        document.getElementById('roleId').value = '';
        document.getElementById('roleName').disabled = false;
        loadAllPermissions();
    }

    modal.classList.add('active');
}

function closeRoleModal() {
    document.getElementById('roleModal').classList.remove('active');
    editingRoleId = null;
}

async function loadAllPermissions() {
    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`${API_URL}/permisos`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (response.ok) {
            const data = await response.json();
            allPermissionsData = data.data || [];
            renderPermissionsGrid([]);
        }
    } catch (error) {
        console.error('Error loading permissions:', error);
    }
}

async function loadPermissionsForRole(roleId) {
    try {
        const token = localStorage.getItem('authToken');
        
        const [permisosRes, roleRes] = await Promise.all([
            fetch(`${API_URL}/permisos`, { headers: { 'Authorization': `Bearer ${token}` } }),
            fetch(`${API_URL}/roles/${roleId}`, { headers: { 'Authorization': `Bearer ${token}` } })
        ]);

        if (permisosRes.ok && roleRes.ok) {
            const permisosData = await permisosRes.json();
            const roleData = await roleRes.json();
            
            allPermissionsData = permisosData.data || [];
            const rolePermisos = roleData.data?.permisos || [];
            renderPermissionsGrid(rolePermisos);
        }
    } catch (error) {
        console.error('Error loading permissions for role:', error);
    }
}

function renderPermissionsGrid(selectedPermisos) {
    const container = document.getElementById('rolePermissions');
    const selectedIds = selectedPermisos.map(p => p.id_permiso);
    
    const modulosMap = {};
    allPermissionsData.forEach(p => {
        if (!modulosMap[p.id_modulo]) {
            modulosMap[p.id_modulo] = {
                nombre: p.modulo_nombre,
                permisos: []
            };
        }
        modulosMap[p.id_modulo].permisos.push(p);
    });

    let html = '';
    Object.entries(modulosMap).forEach(([idModulo, modulo]) => {
        html += `<div class="permission-module-header">${escapeHtml(modulo.nombre)}</div>`;
        modulo.permisos.forEach(p => {
            const checked = selectedIds.includes(p.id_permiso) ? 'checked' : '';
            html += `
                <label class="permission-checkbox">
                    <input type="checkbox" name="permisos" value="${p.id_permiso}" ${checked}>
                    <span>${escapeHtml(p.nombre)}</span>
                </label>
            `;
        });
    });

    container.innerHTML = html;
}

async function handleRoleSubmit(e) {
    e.preventDefault();
    
    // Check permission
    const userPermissions = JSON.parse(localStorage.getItem('userPermissions') || '[]');
    if (!userPermissions.includes('roles_crear') && !userPermissions.includes('roles_editar')) {
        showToast('No tienes permisos para guardar roles', 'error');
        return;
    }
    
    const form = e.target;
    const formData = new FormData(form);
    
    const roleData = {
        nombre: formData.get('nombre'),
        descripcion: formData.get('descripcion'),
        es_rol_sistema: formData.get('es_rol_sistema') === 'on'
    };

    const permisosCheckboxes = form.querySelectorAll('input[name="permisos"]:checked');
    const permisos = Array.from(permisosCheckboxes).map(cb => parseInt(cb.value));

    try {
        const token = localStorage.getItem('authToken');
        const url = editingRoleId 
            ? `${API_URL}/roles/${editingRoleId}`
            : `${API_URL}/roles`;
        
        const method = editingRoleId ? 'PUT' : 'POST';

        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ ...roleData, permisos })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Error al guardar rol');
        }

        showToast(editingRoleId ? 'Rol actualizado exitosamente' : 'Rol creado exitosamente', 'success');
        closeRoleModal();
        loadRoles();
    } catch (error) {
        console.error('Error:', error);
        showToast(error.message, 'error');
    }
}

async function deleteRole(roleId, nombre) {
    // Check permission
    const userPermissions = JSON.parse(localStorage.getItem('userPermissions') || '[]');
    if (!userPermissions.includes('roles_eliminar')) {
        showToast('No tienes permisos para eliminar roles', 'error');
        return;
    }

    if (!confirm(`¿Está seguro de eliminar el rol "${nombre}"?`)) {
        return;
    }

    try {
        const token = localStorage.getItem('authToken');
        
        const response = await fetch(`${API_URL}/roles/${roleId}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Error al eliminar rol');
        }

        showToast('Rol eliminado exitosamente', 'success');
        loadRoles();
    } catch (error) {
        console.error('Error:', error);
        showToast(error.message, 'error');
    }
}

function showToast(message, type = 'success') {
    const container = document.getElementById('toastContainer');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    const icon = type === 'success' 
        ? '<svg class="toast-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>'
        : '<svg class="toast-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>';
    
    toast.innerHTML = `${icon}<span class="toast-message">${escapeHtml(message)}</span>`;
    container.appendChild(toast);

    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transform = 'translateX(100%)';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

window.openRoleModal = openRoleModal;
window.deleteRole = deleteRole;
