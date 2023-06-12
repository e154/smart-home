/* eslint-disable no-console */

import {register} from 'register-service-worker';
import stream from '@/api/stream';
import {UUID} from 'uuid-generator-ts';
import {EventNewWebPushPublicKey} from '@/api/stream_types';

class RegisterServiceWorker {

  private pushNotificationSupported: boolean = false;
  private isRunning: boolean = false;
  private currentID: string = '';

  constructor() {
    const uuid = new UUID();
    this.currentID = uuid.getDashFreeUUID();
    this.pushNotificationSupported = this.isPushNotificationSupported();
  }

  public start() {
    if (!this.pushNotificationSupported) {
      return;
    }
    let owner = this;
    this.initializePushNotifications()
      .then(function(consent) {
        if (consent !== 'granted') {
          return;
        }
        if (owner.isRunning) {
          return;
        }
        owner.isRunning = true;
        stream.subscribe('event_new_webpush_public_key', owner.currentID, owner.onMessage);
        if (localStorage.getItem('vapidPublicKey')) {
          owner.registerServiceWorker();
        } else {
          setTimeout(function() {
            stream.send({id: UUID.createUUID(), query: 'event_get_webpush_public_key'});
          }, 2000);
        }
      });
  }

  public stop() {
    stream.unsubscribe('event_new_webpush_public_key', this.currentID);
  }

  private urlBase64ToUint8Array(base64String: string) {
    const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
    const base64 = (base64String + padding)
      .replace(/\-/g, '+')
      .replace(/_/g, '/');
    const rawData = window.atob(base64);
    return Uint8Array.from([...rawData].map(char => char.charCodeAt(0)));
  }

  private isPushNotificationSupported() {
    return 'serviceWorker' in navigator && 'PushManager' in window;
  }

  private initializePushNotifications() {
    // request user grant to show notification
    return Notification.requestPermission(function(result) {
      return result;
    });
  }

  private subscribeUser(registration: ServiceWorkerRegistration) {
    registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: this.urlBase64ToUint8Array(localStorage.getItem('vapidPublicKey')!),
    })
      .then(function(subscription) {
        // console.log(JSON.stringify(subscription));
        stream.send({
          id: UUID.createUUID(),
          query: 'event_add_webpush_subscription',
          body: btoa(JSON.stringify(subscription))
        });
      })
      .catch(function(error) {
        console.error('Service Worker Error', error);
      });
  }

  private registerServiceWorker() {
    let owner = this;
    register(`${process.env.BASE_URL}service-worker.js`, {
      ready() {
        console.log(
          'App is being served from cache by a service worker.\n' +
          'For more details, visit https://goo.gl/AFskqB'
        );
      },
      registered(registration) {
        console.log('Service worker has been registered.');
        // Routinely check for app updates by testing for a new service worker.
        if (registration.active) {
          owner.subscribeUser(registration);
        } else {
          setTimeout(function() {
            if (registration.active) {
              owner.subscribeUser(registration);
            }
          }, 2000);
        }

        setInterval(() => {
          registration.update();
        }, 1000 * 60 * 60); // hourly checks
      },
      cached() {
        console.log('Content has been cached for offline use.');
      },
      updatefound() {
        console.log('New content is downloading.');
      },
      updated(registration) {
        console.log('New content is available; please refresh.');
        // Add a custom event and dispatch it.
        // Used to display of a 'refresh' banner following a service worker update.
        // Set the event payload to the service worker registration object.
        document.dispatchEvent(
          new CustomEvent('swUpdated', {detail: registration})
        );
      },
      offline() {
        console.log('No internet connection found. App is running in offline mode.');
      },
      error(error) {
        console.error('Error during service worker registration:', error);
      }
    });
  }

  private onMessage(event: EventNewWebPushPublicKey) {
    localStorage.vapidPublicKey = event.public_key;
    console.log(`webpush public key was updated`, event.public_key);
    this.registerServiceWorker();
  }

}// RegisterServiceWorker

const registerServiceWorker = new RegisterServiceWorker();
export default registerServiceWorker;
