
---
title: "Memory"
linkTitle: "memory"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/ram1.png" width="300" >}}

&nbsp;

&nbsp;

В системе **Smart Home** имеется плагин "memory", который позволяет отображать параметры оперативной памяти (RAM). Этот плагин предоставляет доступ к следующим параметрам:

1. `total`: Параметр `total` отображает общий объем оперативной памяти системы в байтах.

2. `free`: Параметр `free` указывает на свободную оперативную память в байтах.

3. `used_percent`: Параметр `used_percent` показывает процент использования оперативной памяти.

Пример использования плагина "memory" для получения параметров оперативной памяти:

```javascript
const ramParams = EntityGetAttributes('memory.memory')
console.log(ramParams.total);
console.log(ramParams.free);
console.log(ramParams.used_percent);
```

Вы можете использовать эти параметры для отображения информации об оперативной памяти в вашем проекте **Smart Home**.
