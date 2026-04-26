const CACHE_NAME = 'comigo-reader-pwa-v1'

const APP_SHELL = [
  '/reader',
  '/script/main.js',
  '/script/styles.css',
  '/script/reader.js',
  '/script/reader_pwa.js',
  '/script/flip_modules/pagination_utils.js',
  '/script/flip_modules/interaction_utils.js',
  '/script/wasm/wasm_exec.js',
  '/script/wasm/archive.wasm',
  '/images/manifest.webmanifest',
  '/images/favicon.png',
  '/images/favicon.ico',
  '/images/pwa-192.png',
  '/images/pwa-512.png',
  '/images/pwa-maskable-512.png',
]

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => cache.addAll(APP_SHELL))
      .then(() => self.skipWaiting()),
  )
})

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys()
      .then((keys) => Promise.all(keys.filter((key) => key !== CACHE_NAME).map((key) => caches.delete(key))))
      .then(() => self.clients.claim()),
  )
})

self.addEventListener('fetch', (event) => {
  const request = event.request
  if (request.method !== 'GET') return

  const url = new URL(request.url)
  if (url.origin !== self.location.origin) return

  if (request.mode === 'navigate' || url.pathname === '/reader') {
    event.respondWith(networkFirst(request, '/reader'))
    return
  }

  if (url.pathname === '/script/reader.js' || url.pathname === '/script/reader_pwa.js' || url.pathname === '/reader-sw.js') {
    event.respondWith(networkFirst(request))
    return
  }

  if (
    url.pathname.startsWith('/script/') ||
    url.pathname.startsWith('/images/')
  ) {
    event.respondWith(cacheFirst(request))
  }
})

async function networkFirst(request, fallbackPath) {
  const cache = await caches.open(CACHE_NAME)
  try {
    const response = await fetch(request)
    if (response.ok) {
      await cache.put(request, response.clone())
      if (fallbackPath) {
        await cache.put(fallbackPath, response.clone())
      }
    }
    return response
  } catch (_) {
    return (await cache.match(request)) || (fallbackPath ? cache.match(fallbackPath) : undefined)
  }
}

async function cacheFirst(request) {
  const cache = await caches.open(CACHE_NAME)
  const cached = await cache.match(request)
  if (cached) return cached

  const response = await fetch(request)
  if (response.ok) {
    await cache.put(request, response.clone())
  }
  return response
}
