import { ROLE_COLORS } from './config.js';

export function showToast(message, type = 'success') {
    const container = document.getElementById('toastContainer') || createToastContainer();
    
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    const iconSvg = type === 'success' 
        ? '<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/>'
        : '<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>';
    
    toast.innerHTML = `
        <svg class="toast-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            ${iconSvg}
        </svg>
        <span class="toast-message">${escapeHtml(message)}</span>
    `;
    
    container.appendChild(toast);
    
    setTimeout(() => {
        toast.style.animation = 'slideOut 0.3s ease forwards';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
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

export function createBadge(text, type = 'default') {
    return `<span class="badge badge-${type}">${escapeHtml(text)}</span>`;
}

export function createRoleBadge(rol) {
    const colorClass = ROLE_COLORS[rol] || 'default';
    return `<span class="role-badge ${colorClass}">${escapeHtml(rol)}</span>`;
}

export function createStatusBadge(activo) {
    const text = activo ? 'Activo' : 'Inactivo';
    const className = activo ? 'active' : 'inactive';
    return `<span class="status-badge ${className}">${text}</span>`;
}

export function createSystemBadge(esSistema) {
    const text = esSistema ? 'Sí' : 'No';
    const className = esSistema ? 'yes' : 'no';
    return `<span class="system-badge ${className}">${text}</span>`;
}

export function createLoadingCell(cols = 1) {
    return `<tr><td colspan="${cols}" class="loading-cell">Cargando...</td></tr>`;
}

export function createEmptyCell(message, cols = 1) {
    return `<tr><td colspan="${cols}" class="empty-cell">${escapeHtml(message)}</td></tr>`;
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
    button.dataset.originalText = button.innerHTML;
    button.disabled = true;
    button.classList.add('loading');
}

export function hideLoading(button) {
    if (!button) return;
    button.disabled = false;
    button.classList.remove('loading');
    if (button.dataset.originalText) {
        button.innerHTML = button.dataset.originalText;
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
