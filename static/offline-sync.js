// Offline sync: queues form submissions made while offline and replays them
// when the device reconnects. Only forms with data-offline-sync are queued.

const DB_NAME = 'mygympal-offline';
const STORE_NAME = 'pending';

function openDB() {
    return new Promise((resolve, reject) => {
        const req = indexedDB.open(DB_NAME, 1);
        req.onupgradeneeded = () => {
            req.result.createObjectStore(STORE_NAME, { keyPath: 'id', autoIncrement: true });
        };
        req.onsuccess = () => resolve(req.result);
        req.onerror = () => reject(req.error);
    });
}

async function savePending(url, formData) {
    const db = await openDB();
    // FormData can hold multiple values per key; capture all of them.
    const data = {};
    for (const [k, v] of formData.entries()) {
        data[k] = v;
    }
    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readwrite');
        tx.objectStore(STORE_NAME).add({ url, data, timestamp: Date.now() });
        tx.oncomplete = resolve;
        tx.onerror = () => reject(tx.error);
    });
}

async function getPending() {
    const db = await openDB();
    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readonly');
        const req = tx.objectStore(STORE_NAME).getAll();
        req.onsuccess = () => resolve(req.result);
        req.onerror = () => reject(req.error);
    });
}

async function clearPending() {
    const db = await openDB();
    return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readwrite');
        tx.objectStore(STORE_NAME).clear();
        tx.oncomplete = resolve;
        tx.onerror = () => reject(tx.error);
    });
}

async function replayPending() {
    const items = await getPending();
    if (!items.length) return;

    let allSucceeded = true;
    for (const item of items) {
        const body = new URLSearchParams(item.data);
        try {
            const res = await fetch(item.url, {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body,
            });
            if (!res.ok) {
                allSucceeded = false;
            }
        } catch {
            // Still offline — stop trying.
            return;
        }
    }

    if (allSucceeded) {
        await clearPending();
        window.location.reload();
    }
}

function showOfflineBanner() {
    if (document.getElementById('offline-banner')) return;
    const el = document.createElement('div');
    el.id = 'offline-banner';
    el.className = 'alert alert-warning alert-dismissible fade show position-fixed top-0 start-50 translate-middle-x mt-3 shadow';
    el.style.zIndex = '9999';
    el.style.minWidth = '300px';
    el.innerHTML =
        "<strong>You're offline.</strong> Your changes have been saved and will sync automatically when you reconnect." +
        '<button type="button" class="btn-close" data-bs-dismiss="alert"></button>';
    document.body.prepend(el);
}

function initFormInterception() {
    // Only intercept forms explicitly marked for offline sync.
    document.querySelectorAll('form[data-offline-sync]').forEach(form => {
        form.addEventListener('submit', async e => {
            if (navigator.onLine) return; // Let the normal submission proceed.
            if (!form.checkValidity()) return; // Let browser validation show errors.
            e.preventDefault();
            await savePending(form.action || window.location.pathname, new FormData(form));
            showOfflineBanner();
        });
    });
}

// Replay any queued submissions on page load (if we're online).
if (navigator.onLine) {
    replayPending();
}

// Also replay when the connection is restored mid-session.
window.addEventListener('online', replayPending);

// Wire up form interception once the DOM is ready.
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initFormInterception);
} else {
    initFormInterception();
}
