// This is the code piece that GenerateSW mode can't provide for us.
// This code listens for the user's confirmation to update the app.
self.addEventListener('install', e => e.waitUntil(self.skipWaiting()));
self.addEventListener('activate', e => e.waitUntil(self.clients.claim()));
self.addEventListener('push', (e) => {
  console.log('Received a push message')
  self.registration.pushManager.getSubscription().then(function (subscription) {
      try {
        let data = e.data.json();
        let options = {
          body: data.body,
        };
        self.registration.showNotification(data.title, options)
      } catch (err) {
        self.registration.showNotification(e.data.text())
      }
    }
  )
});
/* eslint-disable no-undef */
workbox.core.clientsClaim()
workbox.precaching.precacheAndRoute(self.__precacheManifest || [])
