
---
title: "HDD"
linkTitle: "hdd"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/hdd1.png" >}}

&nbsp;

&nbsp;

В системе **Smart Home** присутствует плагин "hdd", который обеспечивает отображение параметров жесткого диска. Этот плагин позволяет получать информацию о различных характеристиках жесткого диска. Вот некоторые из этих параметров:

1. `path`: Параметр `path` содержит путь к точке монтирования (mount point) жесткого диска.

2. `fstype`: Параметр `fstype` указывает тип файловой системы, используемой на жестком диске.

3. `total`: Параметр `total` отображает общий объем жесткого диска в байтах.

4. `free`: Параметр `free` показывает свободное пространство на жестком диске в байтах.

5. `used`: Параметр `used` указывает использованное пространство на жестком диске в байтах.

6. `used_percent`: Параметр `used_percent` показывает процент использования пространства на жестком диске.

7. `inodes_total`: Параметр `inodes_total` отображает общее количество инодов (inode) на жестком диске.

8. `inodes_used`: Параметр `inodes_used` указывает количество использованных инодов на жестком диске.

9. `inodes_free`: Параметр `inodes_free` показывает количество свободных инодов на жестком диске.

10. `inodes_used_percent`: Параметр `inodes_used_percent` указывает процент использования инодов на жестком диске.

Кроме того, плагин "hdd" имеет опцию настроек `mount_point`, которая позволяет указать точку монтирования для отображения параметров конкретного жесткого диска.

Пример использования плагина "hdd" для получения параметров жесткого диска:

```javascript
const hddParams = EntityGetAttributes('hdd.hdd1')
const hddSettings = entity.getSettings()
console.log(hddSettings.mount_point);
console.log(hddParams.path);
console.log(hddParams.fstype);
console.log(hddParams.total);
console.log(hddParams.free);
console.log(hddParams.used);
console.log(hddParams.used_percent);
console.log(hddParams.inodes_total);
console.log(hddParams.inodes_used);
console.log(hddParams.inodes_free);
console.log(hddParams.inodes_used_percent);
```

Вы можете использовать эти параметры для отображения информации о жестком диске в вашем проекте **Smart Home**.
