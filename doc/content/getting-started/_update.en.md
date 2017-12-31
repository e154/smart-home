---
weight: 40
title: overview
groups:
    - "getting_started"
---

<h2 id="update">Обновление</h2>

Обновление системных компонентов заключается в полной замене всех исполняемых файлов, и возмножно частичной
коррекции конфигурационных файлов.

Правило замены справедливо для следующих компонентов:

* Сервер
* Конфигуратор
* Нода

<h2 id="migrate-mysql">Обновление базы mysql</h2>

База поддерживает версионирование, перед обновление структуры нужно создать резервную копию

```bash
./opt/smart-home/server/examples/scripts/backup.sh
Starting backup to 19-02-2017-22-06-30-dump.sql...
```

Скрипт создаст копию действующей базы в */opt/smart-home/backup*

Запуск обновления базы:

```bash
cd /opt/smart-home/server/
./server-linux-amd64 migrate
2017/02/19 01:59:37 GOPATH: /home/delta54/workspace/golang/src/github.com/e154/smart-home/migrate.go 81 /home/delta54/workspace/golang
2017/02/19 01:59:37 Using 'mysql' as 'driver'
2017/02/19 01:59:37 Using 'smarthome:smarthome@tcp(127.0.0.1:3306)/smarthome_dev?charset=utf8' as 'conn'
2017/02/19 01:59:37 Running all outstanding migrations
2017/02/19 01:59:40 |> 2017/02/19 01:59:38 [I] total success upgrade: 0  migration
2017/02/19 01:59:40 Migration successful!
```

Сообщение *Migration successful* говорит о успешно выполненной миграции