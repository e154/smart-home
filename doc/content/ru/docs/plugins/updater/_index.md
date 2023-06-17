
---
title: "Updater"
linkTitle: "updater"
date: 2021-10-20
description: >
  
---

Плагин "updated" позволяет системе проверять наличие обновлений и предоставлять информацию о самой свежей версии.
Это полезно для пользователей, которые хотят быть в курсе последних обновлений и внедрять новые функции или исправления ошибок.
Он отображает следующие параметры:

- `latest_version`: Номер последней доступной версии системы **Smart Home**.
- `latest_version_time`: Дата и время выпуска последней версии.
- `latest_download_url`: URL-адрес, по которому можно скачать последнюю версию системы.
- `last_check`: Дата и время последней проверки наличия обновлений.

Пример использования плагина "updated":

```javascript
const entity = entityManager.getEntity('updater.updater')
const updateInfo = entity.getAttributes()

console.log("Latest version:", updateInfo.latest_version);
console.log("Latest version time:", updateInfo.latest_version_time);
console.log("Latest download URL:", updateInfo.latest_download_url);
console.log("Last check:", updateInfo.last_check);
```

