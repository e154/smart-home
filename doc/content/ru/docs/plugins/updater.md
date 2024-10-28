---
title: "Updater"
linkTitle: "updater"
date: 2021-10-20
description: >
  
---

Плагин "updated" позволяет системе проверять наличие обновлений и предоставлять информацию о самой свежей версии.
Это полезно для пользователей, которые хотят быть в курсе последних обновлений и внедрять новые функции или исправления
ошибок.
Он отображает следующие параметры:

1. `latest_version`: Номер последней доступной версии системы **Smart Home**.
2. `latest_version_time`: Дата и время выпуска последней версии.
3. `latest_download_url`: URL-адрес, по которому можно скачать последнюю версию системы.
4. `last_check`: Дата и время последней проверки наличия обновлений.

Пример использования плагина "updated":

```javascript
const updateInfo = EntityGetAttributes('updater.updater')

console.log("Latest version:", updateInfo.latest_version);
console.log("Latest version time:", updateInfo.latest_version_time);
console.log("Latest download URL:", updateInfo.latest_download_url);
console.log("Last check:", updateInfo.last_check);
```

