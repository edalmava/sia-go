import { ROLE_COLORS } from './config.js';

export function showToast(message, type = 'success') {
    const container = document.getElementById('toastContainer') || createToastContainer();
    
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    
    // Iconos modernos (Lucide-like)
    const icons = {
        success: '<circle cx="12" cy="12" r="10"/><path d="m9 12 2 2 4-4"/>',
        error: '<circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/>',
        info: '<circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/>',
        warning: '<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/>'
    };
    
    toast.innerHTML = `
        <svg class="toast-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            ${icons[type] || icons.info}
        </svg>
        <span class="toast-message">${escapeHtml(message)}</span>
    `;
    
    container.appendChild(toast);
    
    // Auto-remove con animación
    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transform = 'translateX(20px)';
        toast.style.transition = 'all 0.3s ease';
        setTimeout(() => toast.remove(), 300);
    }, 4000);
}

function createToastContainer() {
    const container = document.createElement('div');
    container.id = 'toastContainer';
    container.className = 'toast-container';
    document.body.appendChild(container);
    return container;
}

export function escapeHtml(text) {
    if (text === null || text === undefined) return '';
    const div = document.createElement('div');
    div.textContent = String(text);
    return div.innerHTML;
}

export function formatDate(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('es-CO', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

export function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

export function createBadge(text, type = 'primary') {
    return `<span class="badge badge-${type}">${escapeHtml(text)}</span>`;
}

export function createRoleBadge(rol) {
    const roleMap = {
        'ADMIN': 'primary',
        'DIRECTOR': 'info',
        'DOCENTE': 'warning',
        'ESTUDIANTE': 'success',
        'ACUDIENTE': 'secondary'
    };
    const type = roleMap[rol] || 'primary';
    return `<span class="badge badge-${type}">${escapeHtml(rol)}</span>`;
}

export function createStatusBadge(activo) {
    const type = activo ? 'success' : 'error';
    const text = activo ? 'Activo' : 'Inactivo';
    return `<span class="badge badge-${type}">${text}</span>`;
}

export function createSystemBadge(esSistema) {
    const text = esSistema ? 'Sí' : 'No';
    const type = esSistema ? 'primary' : 'secondary';
    return `<span class="badge badge-${type}">${text}</span>`;
}

export function createLoadingCell(cols = 1) {
    return `<tr><td colspan="${cols}" class="loading-cell">
        <div class="flex items-center justify-center gap-2">
            <svg class="animate-spin" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" style="animation: spin 1s linear infinite;">
                <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
            </svg>
            <span>Cargando datos...</span>
        </div>
    </td></tr>`;
}

export function createEmptyCell(message, cols = 1) {
    return `<tr><td colspan="${cols}" class="empty-cell text-center text-muted">${escapeHtml(message)}</td></tr>`;
}

export function createPagination(total, offset, limit) {
    return {
        total,
        offset,
        limit,
        totalPages: Math.ceil(total / limit),
        currentPage: Math.floor(offset / limit) + 1
    };
}

export function showLoading(button) {
    if (!button) return;
    button.disabled = true;
    button.dataset.originalContent = button.innerHTML;
    button.innerHTML = `
        <svg class="animate-spin" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" style="animation: spin 1s linear infinite;">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
        </svg>
        <span>Cargando...</span>
    `;
}

export function hideLoading(button) {
    if (!button) return;
    button.disabled = false;
    if (button.dataset.originalContent) {
        button.innerHTML = button.dataset.originalContent;
    }
}

export function initModalTriggers(modalId, openBtns, closeBtn, cancelBtn) {
    const modal = document.getElementById(modalId);
    if (!modal) return;

    openBtns.forEach(btn => {
        if (btn) btn.addEventListener('click', () => openModal(modal));
    });

    if (closeBtn) {
        closeBtn.addEventListener('click', () => closeModal(modal));
    }

    if (cancelBtn) {
        cancelBtn.addEventListener('click', () => closeModal(modal));
    }

    modal.addEventListener('click', (e) => {
        if (e.target === modal) closeModal(modal);
    });
}

export function openModal(modal) {
    modal.classList.add('active');
    document.body.style.overflow = 'hidden';
}

export function closeModal(modal) {
    modal.classList.remove('active');
    document.body.style.overflow = '';
}

export function confirmAction(message, onConfirm) {
    if (confirm(message)) {
        onConfirm();
    }
}

export function handleApiError(error, defaultMessage = 'Ha ocurrido un error') {
    console.error('API Error:', error);
    const message = error.data?.message || error.message || defaultMessage;
    showToast(message, 'error');
}
