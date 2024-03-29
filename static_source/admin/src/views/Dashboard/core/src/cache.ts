export class Cache {
  private pull: object;

  constructor() {
    this.pull = {}
  }

  push(key: string, value: any) {
    this.pull[key] = value
  }

  exist(key: string): boolean {
    return this.pull.hasOwnProperty(key)
  }

  get(key: string): any | null {
    if (!this.pull.hasOwnProperty(key)) {
      return null
    }
    return this.pull[key]
  }

  clear() {
    this.pull = {}
  }
}
