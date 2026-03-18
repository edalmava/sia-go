document.addEventListener('DOMContentLoaded', function() {
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

    // Check if user has permission to view users
    const permisos = JSON.parse(localStorage.getItem('userPermissions') || '[]');
    
    if (!permisos.includes('usuarios_ver')) {
        showToast('No tienes permisos para acceder a esta página', 'error');
        setTimeout(() => {
            window.location.href = 'dashboard.html';
        }, 1500);
        return;
    }

    // Show menu items based on permissions
    document.querySelectorAll('.nav-item[data-permiso]').forEach(item => {
        const permiso = item.dataset.permiso;
        if (permisos.includes(permiso)) {
            item.style.display = 'flex';
        }
    });

    const username = localStorage.getItem('username') || 'admin';
    const userRole = localStorage.getItem('userRole') || 'USER';
    
    document.getElementById('userName').textContent = username.charAt(0).toUpperCase() + username.slice(1);
    document.getElementById('userInitial').textContent = username.charAt(0).toUpperCase();
    
    // Set role text
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

    // Show/hide "Nuevo Usuario" button based on permission
    const addUserBtn = document.getElementById('addUserBtn');
    if (addUserBtn) {
        if (permisos.includes('usuarios_crear')) {
            addUserBtn.style.display = 'flex';
        } else {
            addUserBtn.style.display = 'none';
        }
    }

    let users = [];
    let editingUserId = null;
    let deletingUserId = null;

    document.getElementById('logoutBtn').addEventListener('click', function(e) {
        e.preventDefault();
        localStorage.removeItem('authToken');
        localStorage.removeItem('username');
        localStorage.removeItem('userRole');
        localStorage.removeItem('userPermissions');
        window.location.href = 'login.html';
    });

    loadUsers();
    loadDocentes();
    loadEstudiantes();

    document.getElementById('addUserBtn').addEventListener('click', () => openModal());
    document.getElementById('closeModal').addEventListener('click', closeModal);
    document.getElementById('cancelBtn').addEventListener('click', closeModal);
    document.getElementById('userModal').addEventListener('click', function(e) {
        if (e.target === this) closeModal();
    });

    document.getElementById('closePasswordModal').addEventListener('click', closePasswordModal);
    document.getElementById('cancelPasswordBtn').addEventListener('click', closePasswordModal);
    document.getElementById('passwordModal').addEventListener('click', function(e) {
        if (e.target === this) closePasswordModal();
    });

    document.getElementById('closeDeleteModal').addEventListener('click', closeDeleteModal);
    document.getElementById('cancelDeleteBtn').addEventListener('click', closeDeleteModal);
    document.getElementById('deleteModal').addEventListener('click', function(e) {
        if (e.target === this) closeDeleteModal();
    });

    document.getElementById('userForm').addEventListener('submit', handleUserSubmit);
    document.getElementById('passwordForm').addEventListener('submit', handlePasswordSubmit);
    document.getElementById('confirmDeleteBtn').addEventListener('click', handleDelete);

    document.getElementById('searchInput').addEventListener('input', filterUsers);
});

const API_URL = '/api/v1';

function hasPermission(permiso) {
    const permisos = JSON.parse(localStorage.getItem('userPermissions') || '[]');
    return permisos.includes(permiso);
}

async function loadUsers() {
    const tbody = document.getElementById('usersTableBody');
    tbody.innerHTML = '<tr><td colspan="6" class="loading-cell">Cargando usuarios...</td></tr>';

    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`${API_URL}/usuarios`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Error al cargar usuarios');
        }

        const data = await response.json();
        users = data.data || [];
        renderUsers(users);
    } catch (error) {
        console.error('Error:', error);
        tbody.innerHTML = '<tr><td colspan="6" class="empty-cell">Error al cargar usuarios</td></tr>';
        showToast('Error al cargar usuarios', 'error');
    }
}

function renderUsers(usersToRender) {
    const tbody = document.getElementById('usersTableBody');
    
    if (usersToRender.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" class="empty-cell">No hay usuarios registrados</td></tr>';
        return;
    }

    const canChangePassword = hasPermission('usuarios_cambiar_clave');
    const canEdit = hasPermission('usuarios_editar');
    const canDelete = hasPermission('usuarios_eliminar');
    
    tbody.innerHTML = usersToRender.map(user => {
        const roleName = getRoleName(user.id_rol);
        const roleClass = roleName.toLowerCase();
        
        let actionButtons = '';
        
        if (canChangePassword) {
            actionButtons += `
                <button class="btn-action edit" onclick="openPasswordModal(${user.id_usuario}, '${escapeHtml(user.nombre_usuario)}')" title="Cambiar contraseña">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                        <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                    </svg>
                </button>`;
        }
        
        if (canEdit) {
            actionButtons += `
                <button class="btn-action ${user.activo ? 'deactivate' : 'activate'}" onclick="toggleUserStatus(${user.id_usuario}, '${escapeHtml(user.nombre_usuario)}', ${user.activo})" title="${user.activo ? 'Desactivar' : 'Activar'} usuario">
                    ${user.activo 
                        ? `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
                           </svg>`
                        : `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                            <polyline points="22 4 12 14.01 9 11.01"/>
                           </svg>`
                    }
                </button>`;
        }
        
        if (canDelete) {
            actionButtons += `
                <button class="btn-action delete" onclick="openDeleteModal(${user.id_usuario}, '${escapeHtml(user.nombre_usuario)}')" title="Eliminar usuario">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    </svg>
                </button>`;
        }
        
        return `
            <tr>
                <td><strong>${escapeHtml(user.nombre_usuario)}</strong></td>
                <td><span class="role-badge ${roleClass}">${roleName}</span></td>
                <td><span class="status-badge ${user.activo ? 'active' : 'inactive'}">${user.activo ? 'Activo' : 'Inactivo'}</span></td>
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

function openModal(user = null) {
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

    modal.classList.add('active');
}

function closeModal() {
    document.getElementById('userModal').classList.remove('active');
    editingUserId = null;
}

function openPasswordModal(userId, username) {
    document.getElementById('passwordUserId').value = userId;
    document.getElementById('deleteUserName').textContent = username;
    document.getElementById('passwordForm').reset();
    document.getElementById('passwordModal').classList.add('active');
}

function closePasswordModal() {
    document.getElementById('passwordModal').classList.remove('active');
}

function openDeleteModal(userId, username) {
    deletingUserId = userId;
    document.getElementById('deleteUserName').textContent = username;
    document.getElementById('deleteModal').classList.add('active');
}

function closeDeleteModal() {
    document.getElementById('deleteModal').classList.remove('active');
    deletingUserId = null;
}

async function handleUserSubmit(e) {
    e.preventDefault();
    
    const form = e.target;
    const formData = new FormData(form);
    
    const userData = {
        nombre_usuario: formData.get('nombre_usuario'),
        id_rol: parseInt(formData.get('id_rol')),
        activo: formData.get('activo') === 'on'
    };

    const idDocente = formData.get('id_docente');
    const idEstudiante = formData.get('id_estudiante');

    if (idDocente) {
        userData.id_docente = parseInt(idDocente);
    }
    if (idEstudiante) {
        userData.id_estudiante = parseInt(idEstudiante);
    }

    if (!editingUserId) {
        userData.clave = formData.get('clave');
        if (!userData.clave) {
            showToast('La contraseña es obligatoria', 'error');
            return;
        }
    }

    try {
        const token = localStorage.getItem('authToken');
        const url = editingUserId 
            ? `${API_URL}/usuarios/${editingUserId}`
            : `${API_URL}/usuarios`;
        
        const method = editingUserId ? 'PUT' : 'POST';

        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(userData)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Error al guardar usuario');
        }

        showToast(editingUserId ? 'Usuario actualizado exitosamente' : 'Usuario creado exitosamente', 'success');
        closeModal();
        loadUsers();
    } catch (error) {
        console.error('Error:', error);
        showToast(error.message, 'error');
    }
}

async function handlePasswordSubmit(e) {
    e.preventDefault();
    
    const userId = document.getElementById('passwordUserId').value;
    const newPassword = document.getElementById('newPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    if (newPassword !== confirmPassword) {
        showToast('Las contraseñas no coinciden', 'error');
        return;
    }

    if (newPassword.length < 8) {
        showToast('La contraseña debe tener al menos 8 caracteres', 'error');
        return;
    }

    try {
        const token = localStorage.getItem('authToken');
        
        const response = await fetch(`${API_URL}/usuarios/${userId}/change-password`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ password: newPassword })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Error al cambiar contraseña');
        }

        showToast('Contraseña cambiada exitosamente', 'success');
        closePasswordModal();
    } catch (error) {
        console.error('Error:', error);
        showToast(error.message, 'error');
    }
}

async function toggleUserStatus(userId, username, currentStatus) {
    const newStatus = !currentStatus;
    const action = newStatus ? 'activar' : 'desactivar';
    
    if (!confirm(`¿Está seguro de ${action} el usuario "${username}"?`)) {
        return;
    }

    try {
        const token = localStorage.getItem('authToken');
        
        const response = await fetch(`${API_URL}/usuarios/${userId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ activo: newStatus })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || `Error al ${action} usuario`);
        }

        showToast(`Usuario ${newStatus ? 'activado' : 'desactivado'} exitosamente`, 'success');
        loadUsers();
    } catch (error) {
        console.error('Error:', error);
        showToast(error.message, 'error');
    }
}

async function handleDelete() {
    if (!deletingUserId) return;

    try {
        const token = localStorage.getItem('authToken');
        
        const response = await fetch(`${API_URL}/usuarios/${deletingUserId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || 'Error al eliminar usuario');
        }

        showToast('Usuario eliminado exitosamente', 'success');
        closeDeleteModal();
        loadUsers();
    } catch (error) {
        console.error('Error:', error);
        showToast(error.message, 'error');
    }
}

async function loadDocentes() {
    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`${API_URL}/docentes`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (response.ok) {
            const data = await response.json();
            const select = document.getElementById('docente');
            const docentes = data.data || [];
            
            docentes.forEach(docente => {
                const option = document.createElement('option');
                option.value = docente.id_docente;
                option.textContent = `${docente.nombres} ${docente.apellidos}`;
                select.appendChild(option);
            });
        }
    } catch (error) {
        console.log('Error loading docentes:', error);
    }
}

async function loadEstudiantes() {
    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`${API_URL}/estudiantes`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (response.ok) {
            const data = await response.json();
            const select = document.getElementById('estudiante');
            const estudiantes = data.data || [];
            
            estudiantes.forEach(estudiante => {
                const option = document.createElement('option');
                option.value = estudiante.id_estudiante;
                option.textContent = `${estudiante.nombres} ${estudiante.apellidos}`;
                select.appendChild(option);
            });
        }
    } catch (error) {
        console.log('Error loading estudiantes:', error);
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

window.openPasswordModal = openPasswordModal;
window.openDeleteModal = openDeleteModal;
window.toggleUserStatus = toggleUserStatus;
