import { checkAuth, hasPermission, hasAnyPermissionPrefix, initUserInfo, initLogout, initSidebar, getUserData } from '../js/auth.js';
import { rolApi, permisoApi, moduloApi, api } from '../js/api.js';
import { showToast, escapeHtml, openModal, closeModal, handleApiError } from '../js/ui.js';

let allPermisos = [];
let allModulos = [];
let editingRoleId = null;

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) {
        return;
    }

    if (!hasAnyPermissionPrefix(['usuarios_', 'roles_', 'permisos_'])) {
        showToast('No tienes permisos para acceder a esta página', 'error');
        setTimeout(() => {
            window.location.href = 'dashboard.html';
        }, 1500);
        return;
    }

    initUserInfo();
    initSidebar();
    initLogout();
    initTabs();
    initUI();
    initEventListeners();

    loadRoles();
    loadPermisos();
    loadModulos();
});

function initUI() {
    const addRoleBtn = document.getElementById('addRoleBtn');
    if (addRoleBtn) {
        addRoleBtn.style.display = hasPermission('roles_crear') ? 'inline-flex' : 'none';
    }
}

function initTabs() {
    const tabButtons = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-content');
    const { permissions } = getUserData();
    
    const hasRolesPerms = hasPermission('roles_ver') || hasPermission('roles_crear') || hasPermission('roles_editar') || hasPermission('roles_eliminar');
    const hasPermisosPerms = hasPermission('permisos_ver');
    
    tabButtons.forEach(btn => {
        const tabName = btn.dataset.tab;
        
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
            document.getElementById(`${btn.dataset.tab}-tab`)?.classList.add('active');
        });
    });

    if (!hasRolesPerms) {
        document.getElementById('roles-tab')?.style.setProperty('display', 'none');
    }
    if (!hasPermisosPerms) {
        document.getElementById('permisos-tab')?.style.setProperty('display', 'none');
    }

    const visibleTabs = Array.from(tabButtons).filter(btn => btn.style.display !== 'none');
    if (visibleTabs.length > 0) {
        visibleTabs[0].classList.add('active');
        const firstTabName = visibleTabs[0].dataset.tab;
        document.getElementById(`${firstTabName}-tab`)?.classList.add('active');
    }
}

function initEventListeners() {
    document.getElementById('addRoleBtn')?.addEventListener('click', () => openRoleModal());
    document.getElementById('closeRoleModal')?.addEventListener('click', closeRoleModal);
    document.getElementById('cancelRoleBtn')?.addEventListener('click', closeRoleModal);
    document.getElementById('roleModal')?.addEventListener('click', (e) => {
        if (e.target.id === 'roleModal') closeRoleModal();
    });
    document.getElementById('roleForm')?.addEventListener('submit', handleRoleSubmit);
    document.getElementById('permisoSearch')?.addEventListener('input', filterPermisos);
    document.getElementById('moduloFilter')?.addEventListener('change', filterPermisos);
}

async function loadRoles() {
    const tbody = document.getElementById('rolesTableBody');
    tbody.innerHTML = '<tr><td colspan="6" class="loading-cell">Cargando roles...</td></tr>';

    try {
        const data = await rolApi.getAll();
        renderRoles(data.data || []);
    } catch (error) {
        console.error('Error:', error);
        tbody.innerHTML = '<tr><td colspan="6" class="empty-cell">Error al cargar roles</td></tr>';
        handleApiError(error, 'Error al cargar roles');
    }
}

function renderRoles(roles) {
    const tbody = document.getElementById('rolesTableBody');
    const canEditRoles = hasPermission('roles_editar');
    const canDeleteRoles = hasPermission('roles_eliminar');
    
    if (roles.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" class="empty-cell">No hay roles registrados</td></tr>';
        return;
    }

    tbody.innerHTML = roles.map(rol => {
        let actions = '';
        
        if (canEditRoles) {
            actions += `
                <button class="btn-action edit" onclick="openRoleModal(${rol.id_rol}, '${escapeHtml(rol.nombre)}', '${escapeHtml(rol.descripcion || '')}', ${rol.es_rol_sistema})" title="Editar rol">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                </button>`;
        }
        
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
        const data = await permisoApi.getAll();
        allPermisos = data.data || [];
        
        await loadModulosForFilter();
        renderPermisos(allPermisos);
    } catch (error) {
        console.error('Error:', error);
        tbody.innerHTML = '<tr><td colspan="5" class="empty-cell">Error al cargar permisos</td></tr>';
        handleApiError(error, 'Error al cargar permisos');
    }
}

async function loadModulosForFilter() {
    try {
        const data = await moduloApi.getAll();
        allModulos = data.data || [];
        
        const select = document.getElementById('moduloFilter');
        allModulos.forEach(modulo => {
            const option = document.createElement('option');
            option.value = modulo.id_modulo;
            option.textContent = modulo.nombre;
            select.appendChild(option);
        });
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
        const data = await moduloApi.getAll();
        renderModulos(data.data || []);
    } catch (error) {
        console.error('Error:', error);
        grid.innerHTML = '<div class="empty-state">Error al cargar módulos</div>';
        handleApiError(error, 'Error al cargar módulos');
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

let allPermissionsData = [];

window.openRoleModal = async function(roleId = null, nombre = '', descripcion = '', esSistema = false) {
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

    openModal(modal);
};

function closeRoleModal() {
    closeModal(document.getElementById('roleModal'));
    editingRoleId = null;
}

async function loadAllPermissions() {
    try {
        const data = await permisoApi.getAll();
        allPermissionsData = data.data || [];
        renderPermissionsGrid([]);
    } catch (error) {
        console.error('Error loading permissions:', error);
    }
}

async function loadPermissionsForRole(roleId) {
    try {
        const [permisosData, roleData] = await Promise.all([
            permisoApi.getAll(),
            rolApi.getById(roleId)
        ]);
        
        allPermissionsData = permisosData.data || [];
        const rolePermisos = roleData.data?.permisos || [];
        renderPermissionsGrid(rolePermisos);
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
    Object.entries(modulosMap).forEach(([, modulo]) => {
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
    
    if (!hasPermission('roles_crear') && !hasPermission('roles_editar')) {
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
        if (editingRoleId) {
            await rolApi.update(editingRoleId, { ...roleData, permisos });
            showToast('Rol actualizado exitosamente', 'success');
        } else {
            await rolApi.create({ ...roleData, permisos });
            showToast('Rol creado exitosamente', 'success');
        }
        closeRoleModal();
        loadRoles();
    } catch (error) {
        handleApiError(error, 'Error al guardar rol');
    }
}

window.deleteRole = async function(roleId, nombre) {
    if (!hasPermission('roles_eliminar')) {
        showToast('No tienes permisos para eliminar roles', 'error');
        return;
    }

    if (!confirm(`¿Está seguro de eliminar el rol "${nombre}"?`)) {
        return;
    }

    try {
        await rolApi.delete(roleId);
        showToast('Rol eliminado exitosamente', 'success');
        loadRoles();
    } catch (error) {
        handleApiError(error, 'Error al eliminar rol');
    }
};
