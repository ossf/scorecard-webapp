// Minimal service worker to satisfy PWA installability requirements
// (Lighthouse's `installable-manifest`/`service-worker` audits need an
// active service worker controlling the page). Deliberately does not
// implement an offline caching strategy: existing HTTP caching already
// handles static assets, and that's a separate feature decision.

self.addEventListener('install', () => {
  self.skipWaiting()
})

self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim())
})

self.addEventListener('fetch', (event) => {
  event.respondWith(fetch(event.request))
})
