const CACHE_NAME = 'mygympal-v1';

// Pre-cache the offline fallback and CDN assets on install.
const PRECACHE = [
    '/offline',
    'https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css',
    'https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js',
];

self.addEventListener('install', event => {
    event.waitUntil(
        caches.open(CACHE_NAME).then(cache => cache.addAll(PRECACHE))
    );
    self.skipWaiting();
});

self.addEventListener('activate', event => {
    event.waitUntil(
        caches.keys().then(keys =>
            Promise.all(keys.filter(k => k !== CACHE_NAME).map(k => caches.delete(k)))
        )
    );
    self.clients.claim();
});

self.addEventListener('fetch', event => {
    const { request } = event;

    // Only handle GET requests; POST and others pass through normally.
    if (request.method !== 'GET') return;

    const url = new URL(request.url);

    // CDN assets are stable/versioned — serve from cache, populate on first fetch.
    if (url.hostname === 'cdn.jsdelivr.net') {
        event.respondWith(cacheFirst(request));
        return;
    }

    // Same-origin app pages: try network first, fall back to cache, then offline page.
    if (url.origin === self.location.origin) {
        event.respondWith(networkFirst(request));
    }
});

async function cacheFirst(request) {
    const cached = await caches.match(request);
    if (cached) return cached;
    const response = await fetch(request);
    const cache = await caches.open(CACHE_NAME);
    cache.put(request, response.clone());
    return response;
}

async function networkFirst(request) {
    const cache = await caches.open(CACHE_NAME);
    try {
        const response = await fetch(request);
        if (response.ok) {
            cache.put(request, response.clone());
        }
        return response;
    } catch {
        const cached = await cache.match(request);
        if (cached) return cached;
        const offline = await cache.match('/offline');
        return offline || new Response('Offline', { status: 503 });
    }
}
