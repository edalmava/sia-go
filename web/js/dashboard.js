import { checkAuth, initUserInfo, initLogout, initSidebar } from '../js/auth.js';
import { institucionApi } from '../js/api.js';
import { showToast } from '../js/ui.js';

document.addEventListener('DOMContentLoaded', () => {
    if (!checkAuth()) {
        return;
    }

    initUserInfo();
    initSidebar();
    initLogout();
    initNavItems();

    loadDashboardData();
});

function initNavItems() {
    const navItems = document.querySelectorAll('.nav-item');
    navItems.forEach(item => {
        item.addEventListener('click', (e) => {
            const href = item.getAttribute('href');
            if (href && href !== '#') {
                navItems.forEach(nav => nav.classList.remove('active'));
                item.classList.add('active');
                return;
            }
            e.preventDefault();
            navItems.forEach(nav => nav.classList.remove('active'));
            item.classList.add('active');
        });
    });
}

async function loadDashboardData() {
    try {
        const data = await institucionApi.getAll(0, 10);
        console.log('Dashboard data loaded:', data);
    } catch (error) {
        console.log('Using mock data for dashboard');
    }
}
