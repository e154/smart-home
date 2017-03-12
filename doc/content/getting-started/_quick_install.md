---
weight: 25
title: quick install
groups:
    - "getting_started"
---

<h2 id="quick-install" class="page-header">Быстрая установка</h2>

* Подразумевается что у это первая установка в систему
* По командам Ниже
* Установка будет проходить в автоматическом режиме
* Директория установки /opt/smart-home
* После установки потребуется настроить подключение к серверу баз данных mysql 

Установка сервера

```bash
curl -sSL http://e154.github.io/smart-home/server-installer.sh | bash /dev/stdin --install
```

Установка конфигуратор

```bash
curl -sSL http://e154.github.io/smart-home/configurator-installer.sh | bash /dev/stdin --install
```

Установка узла связи

```bash
curl -sSL http://e154.github.io/smart-home/node-installer.sh | bash /dev/stdin --install
```

Нвстройка сервера баз данных mysql

```bash
mysql -u root -p
CREATE DATABASE smarthome;
CREATE USER 'smarthome'@'localhost' IDENTIFIED BY 'smarthome';
GRANT ALL PRIVILEGES ON smarthome . * TO 'smarthome'@'localhost';
FLUSH PRIVILEGES;
use smarthome
source /opt/smart-home/server/dump.sql
```

Запуск сервера

```bash
/opt/smart-home/server/server
```

сервер запустится на порту **3000**

Запуск конфигуратора

```bash
/opt/smart-home/configurator/configurator
```

консоль конфигуратора будет доступа в браузере по адресу [http://localhost:8080](http://localhost:8080) 

Запуск узла связи с устройствами

```bash
/opt/smart-home/node/node
```

Те же команды, но без привязки к консоли

```bash
/opt/smart-home/server/server > /dev/null 2>&1 &
/opt/smart-home/configurator/configurator > /dev/null 2>&1 &
/opt/smart-home/node/node > /dev/null 2>&1 &
```