const API = '/api';

// Auth helpers
function getToken() { return localStorage.getItem('token'); }
function getUser() { 
  try { return JSON.parse(localStorage.getItem('user') || 'null'); } 
  catch { return null; }
}
function setAuth(token, user) {
  localStorage.setItem('token', token);
  localStorage.setItem('user', JSON.stringify(user));
}
function clearAuth() {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
}
function isLoggedIn() { return !!getToken(); }
function isAdmin() { const u = getUser(); return u && u.role === 'admin'; }

async function api(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) };
  const token = getToken();
  if (token) headers['Authorization'] = `Bearer ${token}`;
  
  const res = await fetch(API + path, { ...options, headers });
  const data = await res.json().catch(() => ({}));
  
  if (!res.ok) {
    const err = new Error(data.error || 'Request failed');
    err.status = res.status;
    err.data = data;
    throw err;
  }
  return data;
}

// Cart count
async function updateCartCount() {
  const el = document.getElementById('cart-count');
  if (!el) return;
  if (!isLoggedIn()) { el.textContent = '0'; return; }
  try {
    const data = await api('/cart');
    el.textContent = data.item_count || 0;
  } catch {
    el.textContent = '0';
  }
}

// Navbar auth state
function updateNavAuth() {
  const authBtn = document.getElementById('auth-btn');
  const userMenu = document.getElementById('user-menu');
  if (!authBtn) return;
  
  if (isLoggedIn()) {
    const user = getUser();
    authBtn.style.display = 'none';
    if (userMenu) {
      userMenu.style.display = 'flex';
      const nameEl = document.getElementById('user-name');
      if (nameEl) nameEl.textContent = user.name.split(' ')[0];
    }
  } else {
    authBtn.style.display = 'inline-flex';
    if (userMenu) userMenu.style.display = 'none';
  }
}

function logout() {
  clearAuth();
  window.location.href = '/';
}

// Format price
function formatPrice(n) {
  return '₹' + Number(n).toLocaleString('en-IN');
}

// Star rating
function stars(rating) {
  const full = Math.floor(rating);
  const half = rating % 1 >= 0.5 ? 1 : 0;
  let s = '★'.repeat(full) + (half ? '½' : '') + '☆'.repeat(5 - full - half);
  return s + ` (${rating})`;
}

// Toast
function showToast(msg, type = 'info') {
  let t = document.getElementById('toast');
  if (!t) {
    t = document.createElement('div');
    t.id = 'toast';
    t.style.cssText = 'position:fixed;bottom:24px;right:24px;padding:12px 20px;border-radius:8px;color:white;font-weight:500;z-index:9999;transition:all 0.3s;opacity:0;';
    document.body.appendChild(t);
  }
  t.style.background = type === 'error' ? '#dc3545' : type === 'success' ? '#28a745' : '#0f3460';
  t.textContent = msg;
  t.style.opacity = '1';
  setTimeout(() => t.style.opacity = '0', 3000);
}

// Init on every page
document.addEventListener('DOMContentLoaded', () => {
  updateNavAuth();
  updateCartCount();
});
