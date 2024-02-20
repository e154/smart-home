import stream from "@/api/stream";
import {UUID} from "uuid-generator-ts";
import {EventNewWebPushPublicKey, EventUserDevices, ServerSubscription,} from "@/api/types";
import {GetFullUrl} from "@/utils/serverId";

const uuid = new UUID()

class PushService {
  private worker: IpostMessage | null = null;
  private pushManager: PushManager | null = null;
  private currentID: string = uuid.getDashFreeUUID();
  private serverSubscriptions: ServerSubscription[] | null = null;
  private publicKey: string | null = null;

  constructor() {
  }

  start() {
    setTimeout(() => {
      stream.subscribe('event_user_devices', this.currentID, this.eventUserDevices())
      stream.subscribe('event_new_webpush_public_key', this.currentID, this.eventNewWebpushPublicKey())
      this.registerWorker()
    }, 1000)
  }

  shutdown() {
    stream.unsubscribe('event_user_device', this.currentID)
    stream.unsubscribe('event_new_webpush_public_key', this.currentID)
  }

  public async checkSubscription() {
    console.debug('Проверка подписок ...')
    if (!this.pushManager) {
      console.debug('pushManager не определен')
      return
    }
    try {
      const currentSubscription = await this.pushManager.getSubscription();
      if (currentSubscription) {
        // Проверяем текущую подписку и принимаем решение обновить или использовать её
        const shouldRenewSubscription = this.checkIfSubscriptionNeedsRenewal(currentSubscription);
        if (shouldRenewSubscription) {
          await this.renewSubscription(currentSubscription);
        } else {
          // Используем текущую подписку
          // console.debug(JSON.stringify(currentSubscription))
          console.debug('Используем текущую подписку на уведомления');
        }
      } else {
        // Создаем новую подписку
        await this.createNewSubscription();
      }
    } catch (error) {
      console.error('Ошибка при обработке подписки на уведомления:', error);
    }

  }

  private checkIfSubscriptionNeedsRenewal(clientSubscription: PushSubscription): boolean {

    if (this.serverSubscriptions) {
      // Проверяем подписку на сервере
      const matchingServerSubscription = this.serverSubscriptions.find(serverSubscription => serverSubscription.endpoint === clientSubscription.endpoint);
      if (!matchingServerSubscription) {
        // Подписка не найдена на сервере, нужно обновить подписку
        return true;
      }
    }

    // Если срок действия подписки указан и он истек, то нужно обновить подписку
    if (clientSubscription.expirationTime && clientSubscription.expirationTime < Date.now()) {
      return true;
    }

    // Возвращаем false, если подписка не требует обновления
    return false;
  }

  // ------------------------------------------
  // push manager
  // ------------------------------------------
  private async renewSubscription(clientSubscription: PushSubscription) {
    console.debug('Обновление подписки')
    if (!this.publicKey || !this.pushManager) {
      console.debug('pushManager или publicKey не определен', this.publicKey, this.pushManager)
      return
    }
    try {
      // Получаем новую подписку от сервиса уведомлений
      await this.createNewSubscription()

      // Отменяем текущую подписку, чтобы избежать дублирования уведомлений
      await clientSubscription.unsubscribe();

      console.debug('Подписка успешно обновлена');
    } catch (error) {
      console.error('Ошибка при обновлении подписки на уведомления:', error);
    }
  }

  private async createNewSubscription() {
    console.debug('Создаем новую подписку')
    if (!this.publicKey || !this.pushManager) {
      console.debug('pushManager или publicKey не определен', this.publicKey, this.pushManager)
      return
    }
    // Создаем новую подписку
    const options = {
      applicationServerKey: this.publicKey,
      userVisibleOnly: true,
    }
    const subscription = await this.pushManager.subscribe(options)
    console.debug(JSON.stringify(subscription))
    stream.send({
      id: UUID.createUUID(),
      query: 'event_add_webpush_subscription',
      body: btoa(JSON.stringify(subscription))
    });
  }

  // ------------------------------------------
  // service worker
  // ------------------------------------------

  // private postMessage(msg: IMessageType) {
  //   if (!this.worker) return;
  //   this.worker.postMessage(msg)
  // }
  //
  // private onmessage(ev: MessageEvent) {
  //   //receive user registration data
  // }

  private get getUrl(): string {
    if (window?.app_settings?.server_version) {
      return GetFullUrl('/public/sw.js')
    } else {
      return '/sw.js'
    }
  }

  private async registerWorker() {
    if (!('serviceWorker' in navigator)) {
      throw new Error('No Service Worker support!')
    }
    if (!('PushManager' in window)) {
      throw new Error('No Push API Support!')
    }

    await this.requestNotificationPermission();

    // navigator.serviceWorker.onmessage = this.onmessage

    navigator.serviceWorker.getRegistration(this.getUrl).then((reg: ServiceWorkerRegistration) => {
      if (reg && reg.active) {
        this.worker = reg.active
        this.pushManager = reg.pushManager;
        this.fetchUserDevices()
        this.fetchPublicKey()
        return
      }
      if (!this.worker) {
        navigator.serviceWorker.register(this.getUrl).then((reg: ServiceWorkerRegistration) => {
          if (reg && reg.active) {
            this.worker = reg.active
            this.pushManager = reg.pushManager;
            this.fetchUserDevices()
            this.fetchPublicKey()
          }
        })
      }
    })
  }

  private requestNotificationPermission = async () => {
    const permission = await window.Notification.requestPermission();
    // value of permission can be 'granted', 'default', 'denied'
    // granted: user has accepted the request
    // default: user has dismissed the notification permission popup by clicking on x
    // denied: user has denied the request.
    if (permission !== 'granted') {
      throw new Error('Permission not granted for Notification');
    }
  }

  // ------------------------------------------
  // events
  // ------------------------------------------
  private fetchPublicKey() {
    stream.send({
      id: UUID.createUUID(),
      query: 'event_get_webpush_public_key',
    });
  }

  private fetchUserDevices() {
    stream.send({
      id: UUID.createUUID(),
      query: 'event_get_user_devices',
    });
  }

  private eventNewWebpushPublicKey() {
    return (data: EventNewWebPushPublicKey) => {
      // console.debug(data.public_key)
      if (data.public_key) {
        this.publicKey = data.public_key;
      }
      this.checkSubscription()
    }
  }

  private eventUserDevices() {
    return (data: EventUserDevices) => {
      // console.debug(data)
      if (data) {
        this.serverSubscriptions = data.subscription;
      }
      this.checkSubscription()
    }
  }
}

const pushService: PushService = new PushService()

export default pushService;

export interface IpostMessage {
  postMessage(message: any, transfer: Transferable[]): void;
}

export interface IMessageType {
  t: string
  data: any
}

