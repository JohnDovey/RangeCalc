/* RangeCalc Standard — offline shell + CDN asset cache */
const CACHE = 'range-calc-standard-v1';
const SHELL = new URL('./RangeCalc.html', self.location).href;
const ASSETS = [
  './RangeCalc.html',
  './range-calc.webmanifest',
  './icon.svg',
  'https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css',
  'https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js',
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE).then(async (cache) => {
      for (const url of ASSETS) {
        try {
          await cache.add(url);
        } catch (e) {
          console.warn('[RangeCalc SW] precache skip:', url, e);
        }
      }
      return self.skipWaiting();
    })
  );
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) =>
        Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k)))
      )
      .then(() => self.clients.claim())
  );
});

self.addEventListener('fetch', (event) => {
  const { request } = event;
  if (request.method !== 'GET') return;

  event.respondWith(
    caches.match(request).then((cached) => {
      if (cached) return cached;
      return fetch(request)
        .then((response) => {
          const copy = response.clone();
          if (response.ok) {
            const url = new URL(request.url);
            const samePath =
              url.pathname.endsWith('RangeCalc.html') ||
              url.pathname.endsWith('range-calc.webmanifest') ||
              url.pathname.endsWith('icon.svg');
            const cdn =
              url.hostname === 'cdn.jsdelivr.net' &&
              url.pathname.includes('/bootstrap@5.3.3/');
            if (samePath || cdn) {
              caches.open(CACHE).then((cache) => cache.put(request, copy));
            }
          }
          return response;
        })
        .catch(() => {
          if (request.mode === 'navigate') {
            return caches.match(SHELL);
          }
          return caches.match(request);
        });
    })
  );
});
