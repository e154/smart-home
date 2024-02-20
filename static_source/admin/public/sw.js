import {ExpirationPlugin} from 'workbox-expiration';
import {cleanupOutdatedCaches, precacheAndRoute} from 'workbox-precaching';
import {registerRoute} from 'workbox-routing';
import {CacheFirst} from 'workbox-strategies';
import {CacheableResponsePlugin} from 'workbox-cacheable-response/CacheableResponsePlugin';

// Register precache routes (static cache)
precacheAndRoute(self.__WB_MANIFEST || []);

// Clean up old cache
cleanupOutdatedCaches();

// Google fonts dynamic cache
registerRoute(
  /^https:\/\/fonts\.googleapis\.com\/.*/i,
  new CacheFirst({
    cacheName: "google-fonts-cache",
    plugins: [
      new ExpirationPlugin({maxEntries: 500, maxAgeSeconds: 5184e3}),
      new CacheableResponsePlugin({statuses: [0, 200]})
    ]
  }), "GET");

// Google fonts dynamic cache
registerRoute(
  /^https:\/\/fonts\.gstatic\.com\/.*/i, new CacheFirst({
    cacheName: "gstatic-fonts-cache",
    plugins: [
      new ExpirationPlugin({maxEntries: 500, maxAgeSeconds: 5184e3}),
      new CacheableResponsePlugin({statuses: [0, 200]})
    ]
  }), "GET");

// Dynamic cache for images from `/upload/`
registerRoute(
  /.*upload.*/, new CacheFirst({
    cacheName: "dynamic-images-cache",
    plugins: [
      new ExpirationPlugin({maxEntries: 500, maxAgeSeconds: 5184e3}),
      new CacheableResponsePlugin({statuses: [0, 200]})
    ]
  }), "GET");

// Install and activate service worker
self.addEventListener('install', () => self.skipWaiting());
self.addEventListener('activate', () => self.clients.claim());

// Receive push notifications
self.addEventListener('push', function (e) {
  if (!(
    self.Notification &&
    self.Notification.permission === 'granted'
  )) {
    console.log('notifications aren\'t supported or permission not granted!')
    return;
  }

  if (e.data) {
    try {
      let message = e.data.json();
      self.registration.showNotification(message.title, {
        body: message.body,
        icon: message.icon,
        actions: message.actions
      });
    } catch (e) {
      console.warn(e.data)
    }
  }
});
