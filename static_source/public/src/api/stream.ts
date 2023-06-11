import {ConstantBackoff, Websocket, WebsocketBuilder} from 'websocket-ts';
import {EventHTML5Notify} from '@/api/stream_types';

class Stream {
  private ws: Websocket | null = null;
  protected subscribers: { [key: string]: any };

  constructor() {
    this.subscribers = {};
  }

  connect(url: string, accessToken: string): void {
    if (this.ws) {
      return;
    }
    url = url.replace('http', 'ws') + '/ws?access_token=' + accessToken;
    this.ws = new WebsocketBuilder(url)
      .onOpen((ws: Websocket, ev: Event) => this.onOpen(ws, ev, accessToken))
      .onClose((ws: Websocket, ev: CloseEvent) => this.onClose(ws, ev))
      .onError((ws: Websocket, ev: Event) => this.error(ws, ev))
      .onMessage((ws: Websocket, ev: MessageEvent) => this.onMessage(ws, ev))
      .withBackoff(new ConstantBackoff(1000)) // 1000ms = 1s
      .onRetry((i, ev) => {
        console.log('retry');
      })
      .build();
  }

  disconnect(): void {
    if (!this.ws) {
      return;
    }
    this.ws.close();
    this.ws = null;
  }

  private error(ws: Websocket, ev: Event): any {
    console.log('error');
  }

  private onClose(ws: Websocket, ev: CloseEvent): any {
    console.log('closed');
  }

  private onOpen(ws: Websocket, ev: Event, accessToken: string): any {
    console.log('opened');
    ws.send(JSON.stringify({body: btoa('init'), access_token: accessToken}));
  }

  private onMessage(ws: Websocket, ev: MessageEvent): any {
    let m: StreamResponse;
    try {
      const {result} = JSON.parse(ev.data);
      m = result;
      const body: any = JSON.parse(atob(m.body));
      m.body = body;
    } catch {
      console.debug('from the stream came a string value');
      return;
    }

    if (m.query == 'html5_notify') {
      const {result} = JSON.parse(ev.data);
      m = result;
      const body: EventHTML5Notify = JSON.parse(atob(m.body));
      // console.log(body)
      this.notify(body);
      return;
    }

    if (!this.subscribers) {
      return;
    }

    if (m.query in this.subscribers) {
      for (const id in this.subscribers[m.query]) {
        this.subscribers[m.query][id](m.body);
      }
    }
  }

  public send(data: any) {
    if (this.ws) {
      this.ws.send(JSON.stringify(data));
    }
  }

  public subscribe(name: string, id: string, f: any) {
    console.debug('subscribe', name, id);
    if (this.subscribers[name] === undefined) {
      this.subscribers[name] = {};
    }
    this.subscribers[name][id] = f;
  }

  public unsubscribe(name: string, id: string) {
    console.debug('unsubscribe', name, id);
    if (this.subscribers[name] !== undefined) {
      if (this.subscribers[name][id] !== undefined) {
        delete this.subscribers[name][id];
      }
    }
  }

  public notify(event: EventHTML5Notify) {
    if (!('Notification' in window)) {
      console.warn('the browser does not support HTML notifications, it needs to be updated');
      return;
    }
    if (Notification.permission === 'granted') {
      new Notification(event.title, event.options);
      return;
    }
    if (Notification.permission !== 'denied') {
      Notification.requestPermission(function(permission) {
        if (permission === 'granted') {
          new Notification(event.title, event.options);
        }
      });
    }
  }
}

const stream = new Stream();

export default stream;

export interface StreamResponse {
  id: string;
  query: string;
  body: string;
}

export interface StreamRequest {
  id: string;
  query: string;
  body: string;
}
