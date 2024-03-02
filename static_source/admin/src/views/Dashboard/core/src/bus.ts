import mitt from 'mitt'
import { onBeforeUnmount } from 'vue'

interface Option {
  name: string // 事件名称
  callback: Fn // 回调
}

const bus = mitt()

export const useBus = (option?: Option) => {
  if (option) {
    bus.on(option.name, option.callback)

    onBeforeUnmount(() => {
      bus.off(option.name)
    })
  }

  return {
    on: bus.on,
    off: bus.off,
    emit: bus.emit,
    all: bus.all
  }
}
