/* RangeCalc GPS — offline app shell + Bootstrap + Leaflet (map tiles still need network) */
const CACHE = 'range-calc-gps-v1';
const SHELL = new URL('./RangeCalcGPS.html', self.location).href;
const ASSETS = [
  './RangeCalcGPS.html',
  './range-calc-gps.webmanifest',
  './icon.svg',
  'https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css',
  'https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js',
  'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css',
  'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js',
  'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE).then(async (cache) => {
      for (const url of ASSETS) {
        try {
          await cache.add(url);
        } catch (e) {
          console.warn('[RangeCalc GPS SW] precache skip:', url, e);
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

function shouldRuntimeCache(url) {
  if (
    url.pathname.endsWith('RangeCalcGPS.html') ||
    url.pathname.endsWith('range-calc-gps.webmanifest') ||
    url.pathname.endsWith('icon.svg')
  ) {
    return true;
  }
  if (url.hostname === 'cdn.jsdelivr.net' && url.pathname.includes('/bootstrap@5.3.3/')) {
    return true;
  }
  if (url.hostname === 'unpkg.com' && url.pathname.startsWith('/leaflet@1.9.4/')) {
    return true;
  }
  return false;
}

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
            if (shouldRuntimeCache(url)) {
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
