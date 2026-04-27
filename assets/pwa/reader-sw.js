const CACHE_NAME = 'comigo-reader-pwa-v1'

const BASE_PATH = (() => {
  const scopePath = new URL(self.registration.scope).pathname
  return scopePath.replace(/\/reader\/?$/, '')
})()

function withBase(path) {
  if (!BASE_PATH) return path
  return path === '/' ? BASE_PATH + '/' : BASE_PATH + path
}

function relativePath(pathname) {
  if (!BASE_PATH) return pathname
  if (pathname === BASE_PATH) return '/'
  return pathname.startsWith(BASE_PATH + '/') ? pathname.slice(BASE_PATH.length) : pathname
}

const APP_SHELL = [
  '/reader',
  '/assets/dist/main.js',
  '/assets/dist/styles.css',
  '/assets/static/js/reader.js',
  '/assets/static/js/reader_pwa.js',
  '/assets/static/js/flip_modules/pagination_utils.js',
  '/assets/static/js/flip_modules/interaction_utils.js',
  '/assets/static/wasm/wasm_exec.js',
  '/assets/static/wasm/archive.wasm',
  '/images/manifest.webmanifest',
  '/images/favicon.png',
  '/images/favicon.ico',
  '/images/pwa-192.png',
  '/images/pwa-512.png',
  '/images/pwa-maskable-512.png',
].map(withBase)

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

  const pathname = relativePath(url.pathname)

  if (request.mode === 'navigate' || pathname === '/reader') {
    event.respondWith(networkFirst(request, withBase('/reader')))
    return
  }

  if (pathname === '/assets/static/js/reader.js' || pathname === '/assets/static/js/reader_pwa.js' || pathname === '/reader-sw.js') {
    event.respondWith(networkFirst(request))
    return
  }

  if (
    pathname.startsWith('/assets/') ||
    pathname.startsWith('/images/')
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
