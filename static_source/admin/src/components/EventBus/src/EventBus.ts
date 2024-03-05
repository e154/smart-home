/*

// Пример использования

// Создаем экземпляр EventBus
const eventBus = new EventBus();

// Подписываемся на события "message" и "notification"
eventBus.subscribe(["message", "notification"], (eventName, ...args) => {
  if (eventName === "message") {
    console.log(`Received event "${eventName}":`, ...args);
  } else if (eventName === "notification") {
    console.log(`Received event "${eventName}" from ${args[1]}: ${args[0]}`);
  }
});

// Генерируем событие "message"
eventBus.emit("message", "Hello, world!");

// Генерируем событие "notification"
eventBus.emit("notification", "New message received", "John");

*/

export type EventHandler = (...args: any[]) => void;

export class EventBus {
    private listeners: { [event: string]: EventHandler[] } = {};

    private subscribeToAll(handler: EventHandler) {
        Object.keys(this.listeners).forEach(event => {
            this.listeners[event].push(handler);
        });
    }

    subscribe(events: string | string[] | undefined, handler: EventHandler) {
        if (!events) {
            this.subscribeToAll(handler);
        } else {
            if (Array.isArray(events)) {
                events.forEach(event => {
                    if (!this.listeners[event]) {
                        this.listeners[event] = [];
                    }
                    this.listeners[event].push(handler);
                });
            } else {
                if (!this.listeners[events]) {
                    this.listeners[events] = [];
                }
                this.listeners[events].push(handler);
            }
        }
    }

    unsubscribe(events: string | string[], handler: EventHandler) {
        if (Array.isArray(events)) {
            events.forEach(event => {
                if (this.listeners[event]) {
                    this.listeners[event] = this.listeners[event].filter(h => h !== handler);
                }
            });
        } else {
            if (this.listeners[events]) {
                this.listeners[events] = this.listeners[events].filter(h => h !== handler);
            }
        }
    }

    emit(event: string, ...args: any[]) {
        const eventListeners = this.listeners[event];
        if (eventListeners) {
            eventListeners.forEach(handler => handler(event, ...args));
        }
    }
}
