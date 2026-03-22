import { checkAuth, hasPermission, initUserInfo, initLogout, initSidebar, getUserData } from '../js/auth.js';
import { usuarioApi, rolApi, api } from '../js/api.js';
import { showToast, escapeHtml, openModal, closeModal, handleApiError, createStatusBadge, createRoleBadge } from '../js/ui.js';

let users = [];
let editingUserId = null;
let deletingUserId = null;

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) {
        return;
    }

    if (!hasPermission('usuarios_ver')) {
        showToast('No tienes permisos para acceder a esta página', 'error');
        setTimeout(() => {
            window.location.href = 'dashboard.html';
        }, 1500);
        return;
    }

    initUserInfo();
    initSidebar();
    initLogout();
    initUI();
    initEventListeners();

    loadUsers();
    loadRoles();
    loadDocentes();
    loadEstudiantes();
});

async function loadRoles() {
    try {
        const data = await rolApi.getAll();
        const select = document.getElementById('role');
        if (!select) return;

        const roles = data.data || [];
        // Limpiamos las opciones actuales (excepto la primera)
        select.innerHTML = '<option value="">Seleccione un rol</option>';
        
        roles.forEach(rol => {
            const option = document.createElement('option');
            option.value = rol.id_rol;
            option.textContent = rol.nombre;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Error loading roles:', error);
    }
}

function initUI() {
    const addUserBtn = document.getElementById('addUserBtn');
    if (addUserBtn && hasPermission('usuarios_crear')) {
        addUserBtn.style.display = 'inline-flex';
    }
}

function initEventListeners() {
    document.getElementById('addUserBtn')?.addEventListener('click', () => openUserModal());
    document.getElementById('closeModal')?.addEventListener('click', closeUserModal);
    document.getElementById('cancelBtn')?.addEventListener('click', closeUserModal);
    document.getElementById('userModal')?.addEventListener('click', (e) => {
        if (e.target.id === 'userModal') closeUserModal();
    });

    document.getElementById('closePasswordModal')?.addEventListener('click', closePasswordModal);
    document.getElementById('cancelPasswordBtn')?.addEventListener('click', closePasswordModal);
    document.getElementById('passwordModal')?.addEventListener('click', (e) => {
        if (e.target.id === 'passwordModal') closePasswordModal();
    });

    document.getElementById('closeDeleteModal')?.addEventListener('click', closeDeleteModal);
    document.getElementById('cancelDeleteBtn')?.addEventListener('click', closeDeleteModal);
    document.getElementById('deleteModal')?.addEventListener('click', (e) => {
        if (e.target.id === 'deleteModal') closeDeleteModal();
    });

    document.getElementById('userForm')?.addEventListener('submit', handleUserSubmit);
    document.getElementById('passwordForm')?.addEventListener('submit', handlePasswordSubmit);
    document.getElementById('confirmDeleteBtn')?.addEventListener('click', handleDelete);
    document.getElementById('searchInput')?.addEventListener('input', filterUsers);
}

async function loadUsers() {
    const tbody = document.getElementById('usersTableBody');
    tbody.innerHTML = '<tr><td colspan="6" class="loading-cell">Cargando usuarios...</td></tr>';

    try {
        const data = await usuarioApi.getAll(0, 100);
        users = data.data || [];
        renderUsers(users);
    } catch (error) {
        console.error('Error:', error);
        tbody.innerHTML = '<tr><td colspan="6" class="empty-cell">Error al cargar usuarios</td></tr>';
        handleApiError(error, 'Error al cargar usuarios');
    }
}

function renderUsers(usersToRender) {
    const tbody = document.getElementById('usersTableBody');
    const canChangePassword = hasPermission('usuarios_cambiar_clave');
    const canEdit = hasPermission('usuarios_editar');
    const canDelete = hasPermission('usuarios_eliminar');
    
    if (usersToRender.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" class="empty-cell">No hay usuarios registrados</td></tr>';
        return;
    }

    tbody.innerHTML = usersToRender.map(user => {
        const roleName = getRoleName(user.id_rol);
        const roleClass = roleName.toLowerCase();
        
        let actionButtons = '';
        
        if (canChangePassword) {
            actionButtons += `
                <button class="btn btn-icon btn-ghost" onclick="openPasswordModal(${user.id_usuario}, '${escapeHtml(user.nombre_usuario)}')" title="Cambiar contraseña">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="18" height="18">
                        <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                        <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                    </svg>
                </button>`;
        }
        
        if (canEdit) {
            actionButtons += `
                <button class="btn btn-icon btn-ghost" onclick="toggleUserStatus(${user.id_usuario}, '${escapeHtml(user.nombre_usuario)}', ${user.activo})" title="${user.activo ? 'Desactivar' : 'Activar'} usuario">
                    ${user.activo 
                        ? `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="18" height="18">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
                           </svg>`
                        : `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="18" height="18">
                            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                            <polyline points="22 4 12 14.01 9 11.01"/>
                           </svg>`
                    }
                </button>`;
        }
        
        if (canDelete) {
            actionButtons += `
                <button class="btn btn-icon btn-ghost" style="color: var(--error);" onclick="openDeleteModal(${user.id_usuario}, '${escapeHtml(user.nombre_usuario)}')" title="Eliminar usuario">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="18" height="18">
                        <path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                    </svg>
                </button>`;
        }
        
        return `
            <tr>
                <td><strong>${escapeHtml(user.nombre_usuario)}</strong></td>
                <td>${createRoleBadge(roleName)}</td>
                <td>${createStatusBadge(user.activo)}</td>
                <td>${user.id_docente ? `#${user.id_docente}` : '-'}</td>
                <td>${user.id_estudiante ? `#${user.id_estudiante}` : '-'}</td>
                <td>
                    <div class="action-buttons">
                        ${actionButtons}
                    </div>
                </td>
            </tr>
        `;
    }).join('');
}

function getRoleName(idRol) {
    const roles = {
        1: 'ADMIN',
        2: 'DIRECTOR',
        3: 'DOCENTE',
        4: 'ESTUDIANTE',
        5: 'ACUDIENTE'
    };
    return roles[idRol] || 'USER';
}

function filterUsers() {
    const searchTerm = document.getElementById('searchInput').value.toLowerCase();
    const filtered = users.filter(user => 
        user.nombre_usuario.toLowerCase().includes(searchTerm) ||
        getRoleName(user.id_rol).toLowerCase().includes(searchTerm)
    );
    renderUsers(filtered);
}

function openUserModal(user = null) {
    const modal = document.getElementById('userModal');
    const form = document.getElementById('userForm');
    const title = document.getElementById('modalTitle');
    const passwordGroup = document.getElementById('passwordGroup');
    const passwordInput = document.getElementById('password');

    editingUserId = null;
    form.reset();

    if (user) {
        editingUserId = user.id_usuario;
        title.textContent = 'Editar Usuario';
        passwordGroup.style.display = 'none';
        passwordInput.required = false;
        
        document.getElementById('username').value = user.nombre_usuario;
        document.getElementById('role').value = user.id_rol;
        document.getElementById('docente').value = user.id_docente || '';
        document.getElementById('estudiante').value = user.id_estudiante || '';
        document.getElementById('activo').checked = user.activo;
    } else {
        title.textContent = 'Nuevo Usuario';
        passwordGroup.style.display = 'block';
        passwordInput.required = true;
    }

    openModal(modal);
}

function closeUserModal() {
    closeModal(document.getElementById('userModal'));
    editingUserId = null;
}

window.openPasswordModal = function(userId, username) {
    document.getElementById('passwordUserId').value = userId;
    openModal(document.getElementById('passwordModal'));
};

function closePasswordModal() {
    closeModal(document.getElementById('passwordModal'));
}

window.openDeleteModal = function(userId, username) {
    deletingUserId = userId;
    document.getElementById('deleteUserName').textContent = username;
    openModal(document.getElementById('deleteModal'));
};

function closeDeleteModal() {
    closeModal(document.getElementById('deleteModal'));
    deletingUserId = null;
}

async function handleUserSubmit(e) {
    e.preventDefault();
    
    const form = e.target;
    const submitBtn = form.querySelector('button[type="submit"]');
    
    if (submitBtn.disabled) return;
    submitBtn.disabled = true;
    
    if (!hasPermission('usuarios_crear') && !hasPermission('usuarios_editar')) {
        showToast('No tienes permisos para guardar usuarios', 'error');
        submitBtn.disabled = false;
        return;
    }
    
    const formData = new FormData(form);
    
    const userData = {
        nombre_usuario: formData.get('nombre_usuario'),
        id_rol: parseInt(formData.get('id_rol')),
        activo: formData.get('activo') === 'on'
    };

    const idDocente = formData.get('id_docente');
    const idEstudiante = formData.get('id_estudiante');

    if (idDocente) userData.id_docente = parseInt(idDocente);
    if (idEstudiante) userData.id_estudiante = parseInt(idEstudiante);

    if (!editingUserId) {
        userData.clave = formData.get('clave');
        if (!userData.clave) {
            showToast('La contraseña es obligatoria', 'error');
            submitBtn.disabled = false;
            return;
        }
    }

    try {
        if (editingUserId) {
            await usuarioApi.update(editingUserId, userData);
            showToast('Usuario actualizado exitosamente', 'success');
        } else {
            await usuarioApi.create(userData);
            showToast('Usuario creado exitosamente', 'success');
        }
        closeUserModal();
        loadUsers();
    } catch (error) {
        handleApiError(error, 'Error al guardar usuario');
    } finally {
        submitBtn.disabled = false;
    }
}

async function handlePasswordSubmit(e) {
    e.preventDefault();
    
    const form = e.target;
    const submitBtn = form.querySelector('button[type="submit"]');
    
    if (submitBtn.disabled) return;
    submitBtn.disabled = true;
    
    if (!hasPermission('usuarios_cambiar_clave')) {
        showToast('No tienes permisos para cambiar contraseñas', 'error');
        submitBtn.disabled = false;
        return;
    }
    
    const userId = document.getElementById('passwordUserId').value;
    const newPassword = document.getElementById('newPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    if (newPassword !== confirmPassword) {
        showToast('Las contraseñas no coinciden', 'error');
        submitBtn.disabled = false;
        return;
    }

    if (newPassword.length < 8) {
        showToast('La contraseña debe tener al menos 8 caracteres', 'error');
        submitBtn.disabled = false;
        return;
    }

    try {
        await usuarioApi.changePassword(parseInt(userId), newPassword);
        showToast('Contraseña cambiada exitosamente', 'success');
        closePasswordModal();
    } catch (error) {
        handleApiError(error, 'Error al cambiar contraseña');
    } finally {
        submitBtn.disabled = false;
    }
}

window.toggleUserStatus = async function(userId, username, currentStatus) {
    if (!hasPermission('usuarios_editar')) {
        showToast('No tienes permisos para editar usuarios', 'error');
        return;
    }
    
    const user = users.find(u => u.id_usuario === userId);
    if (!user) {
        showToast('No se encontró la información del usuario', 'error');
        return;
    }

    const newStatus = !currentStatus;
    const action = newStatus ? 'activar' : 'desactivar';
    
    if (!confirm(`¿Está seguro de ${action} el usuario "${username}"?`)) {
        return;
    }

    try {
        await usuarioApi.update(userId, {
            nombre_usuario: user.nombre_usuario,
            id_rol: user.id_rol,
            id_docente: user.id_docente,
            id_estudiante: user.id_estudiante,
            activo: newStatus
        });
        showToast(`Usuario ${newStatus ? 'activado' : 'desactivado'} exitosamente`, 'success');
        loadUsers();
    } catch (error) {
        handleApiError(error, `Error al ${action} usuario`);
    }
};

async function handleDelete() {
    const deleteBtn = document.getElementById('confirmDeleteBtn');
    if (deleteBtn.disabled) return;
    deleteBtn.disabled = true;
    
    if (!hasPermission('usuarios_eliminar')) {
        showToast('No tienes permisos para eliminar usuarios', 'error');
        deleteBtn.disabled = false;
        return;
    }
    
    if (!deletingUserId) {
        deleteBtn.disabled = false;
        return;
    }

    try {
        await usuarioApi.delete(deletingUserId);
        showToast('Usuario eliminado exitosamente', 'success');
        closeDeleteModal();
        loadUsers();
    } catch (error) {
        handleApiError(error, 'Error al eliminar usuario');
    } finally {
        deleteBtn.disabled = false;
    }
}

async function loadDocentes() {
    try {
        const data = await api.get('/docentes');
        const select = document.getElementById('docente');
        const docentes = data.data || [];
        
        docentes.forEach(docente => {
            const option = document.createElement('option');
            option.value = docente.id_docente;
            option.textContent = `${docente.nombres} ${docente.apellidos}`;
            select.appendChild(option);
        });
    } catch (error) {
        console.log('Error loading docentes:', error);
    }
}

async function loadEstudiantes() {
    try {
        const data = await api.get('/estudiantes');
        const select = document.getElementById('estudiante');
        const estudiantes = data.data || [];
        
        estudiantes.forEach(estudiante => {
            const option = document.createElement('option');
            option.value = estudiante.id_estudiante;
            option.textContent = `${estudiante.nombres} ${estudiante.apellidos}`;
            select.appendChild(option);
        });
    } catch (error) {
        console.log('Error loading estudiantes:', error);
    }
}
