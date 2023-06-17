
---
title: "Version"
linkTitle: "version"
date: 2021-10-20
description: >
  
---

Плагин "version" в системе предоставляет информацию о версии и сборке проекта. Он отображает следующие параметры:

- `version`: Версия проекта.
- `revision`: Идентификатор ревизии (commit) в системе контроля версий.
- `revision_url`: URL-адрес, по которому можно просмотреть детали ревизии.
- `generated`: Дата и время генерации сборки проекта.
- `developers`: Список разработчиков, вовлеченных в проект.
- `build_num`: Номер сборки проекта.
- `docker_image`: Название Docker-образа, связанного с проектом.
- `go_version`: Версия языка программирования Go, используемого для разработки проекта.

Плагин "version" предоставляет доступ к этой информации, что полезно для отслеживания версий и контроля процесса разработки.

Пример использования плагина "version":


```javascript
const entity = entityManager.getEntity('version.version')
const versionInfo = entity.getAttributes()
console.log("Project version:", versionInfo.version);
console.log("Revision:", versionInfo.revision);
console.log("Revision URL:", versionInfo.revision_url);
console.log("Generated:", versionInfo.generated);
console.log("Developers:", versionInfo.developers);
console.log("Build number:", versionInfo.build_num);
console.log("Docker image:", versionInfo.docker_image);
console.log("Go version:", versionInfo.go_version);
```

Этот плагин позволяет получить актуальную информацию о версии проекта и использовать ее для различных целей, таких как отображение версии в пользовательском интерфейсе или логирование.
