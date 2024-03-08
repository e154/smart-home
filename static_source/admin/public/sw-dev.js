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
