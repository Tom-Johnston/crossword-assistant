var CACHE_NAME = 'crossword-assistant-v1.0.1';

self.addEventListener('install', function(event) {
  // Perform install steps
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(function(cache) {
        return cache.addAll([
            "./",
            "index.html",
            "css/styles.css",
            "script.js",
            "crossword-assistant.wasm",
            "worker.js",
            "wasm_exec.js",
            "images/install_icon.svg",
            "images/apple-touch-icon.png",
            "fonts/courier-prime-bold.ttf",
            "fonts/courier-prime-regular.ttf",
            "favicon.ico",
            "manifest.webmanifest",
        ]);
      })
  );
});

self.addEventListener('activate', function(event) {
    event.waitUntil(
      caches.keys().then(function(cacheNames) {
        return Promise.all(
          cacheNames.filter(function(cacheName) {
            return cacheName !== CACHE_NAME
          }).map(function(cacheName) {
            return caches.delete(cacheName);
          })
        );
      })
    );
  });

  self.addEventListener('fetch', function(event) {
    event.respondWith(
      caches.match(event.request, {ignoreSearch: true})
        .then(function(response) {
          if (response) {
            return response;
          }
          return fetch(event.request);
        }
      )
    );
  });
