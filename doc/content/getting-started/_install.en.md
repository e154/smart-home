---
weight: 30
title: overview
groups:
    - "getting_started"
---

<h2 id="install">Установка</h2>

Будет произведена базовая установка системы **Умный дом** на сервер под операционной системой linux Debian.
С не значительными изменениями установка будет подобна для других операционных систем.

План действий:
    
*   <a href="#install-mkdir">Создание директорий</a>
*   <a href="#install-download-and-unpack">Скачивание и распаковка</a>
*   <a href="#install-server-conf">Настройка сервера</a>
*   <a href="#install-configuration-conf">Настройка конфигуратора</a>
*   <a href="#install-node-conf">Настройка ноды</a>
*   <a href="#install-postgres">База postgres</a>
    

<div class="row">
    <div class="col-md-6">
        <img src="/smart-home/img/default_network2.png" alt="Схематичная сеть умного дома" title="Схематичная сеть умного дома">
    </div>
</div>


Информация о системе на которой производилось развёртывание:

```bash
delta54@darkstar:tmp$ uname -a
Linux darkstar 3.16.0-4-amd64 #1 SMP Debian 3.16.36-1+deb8u2 (2016-10-19) x86_64 GNU/Linux
```

```bash
delta54@darkstar:tmp$ cat /etc/*release*
PRETTY_NAME="Debian GNU/Linux 8 (jessie)"
```

<h3 id="install-mkdir">Создание директорий</h3>

Рекомендуемая директория установки сервера */opt/smart-home*, в неё и будет производиться дальнейшая установка. 
Подготовка директорий, выставляются требуемые права доступа для */opt/smart-home*
    

```bash
sudo mkdir -p /opt/smart-home
sudo chown $USER:$USER /opt/smart-home -R
mkdir -p /opt/smart-home/server
mkdir -p /opt/smart-home/configurator
mkdir -p /opt/smart-home/node
```

<h3 id="install-download-and-unpack">Скачивание</h3>

```bash
cd /opt/smart-home/server

wget https://github.com/e154/smart-home/releases/download/.../smart-home-server.tar.gz
tar -zxvf smart-home-server.tar.gz

cd /opt/smart-home/server
mv data ../

cd /opt/smart-home/configurator
wget https://github.com/e154/smart-home-configurator/releases/download/.../smart-home-configurator.tar.gz
tar -zxvf smart-home-configurator.tar.gz

cd /opt/smart-home/node
wget https://github.com/e154/smart-home-node/releases/download/.../smart-home-node.tar.gz
tar -zxvf smart-home-node.tar.gz
```

Корень директории */opt/smart-home*:

```bash
cd /opt/smart-home
tree
├── data
├── node
├── configurator
└── server
```

<h3 id="install-server-conf">Настройка сервера</h3>

```bash
cd /opt/smart-home/server
cp conf/config.dev.conf conf/config.conf
```

Важные поля:
    
*   **server_host** - адрес сервера REST API
*   **server_port** - порт сервера REST API
*   **pg_user** - пользователь для соединени postgresql
*   **pg_pass** - пароль доступка к postgresql
*   **pg_host** - адрес сервера postgresql
*   **pg_name** - название базы postgresql
*   **pg_port** - порт сервера postgresql    

```bash
pg_user = smart_home
pg_pass = smart_home
pg_host = "127.0.0.1"
pg_name = smart_home
pg_port = "5432"
```

<h3 id="install-configuration-conf">Настройка конфигуратора</h3>

```bash
cd /opt/smart-home/configurator
cp conf/config.dev.conf conf/config.conf
```

Основные значения:
    
*   **httpaddr** - адрес веб интерфеса конфигуратора
*   **httpport** - порт веб интерфеса конфигуратора
*   **api_addr** - порт сервера REST API
*   **api_port** - порт сервера REST API
*   **api_scheme** - схема общения с сервером
    

<h3 id="install-node-conf">Настройка ноды</h3>

```bash
cd /opt/smart-home/node
cp conf/config.dev.conf conf/config.conf
```

Основные значения:
    
*   **name** - системное назание ноды (прим. node1)
*   **topic** - канал для общения с сервером (прим. node1)
*   **mqtt_keep_alive** - тонкие настройки (прим. 300)
*   **mqtt_connect_timeout** - тонкие настройки (прим. 2)
*   **mqtt_sessions_provider** - тонкие настройки (прим. "mem")
*   **mqtt_topics_provider** - тонкие настройки (прим. "mem)
*   **mqtt_username** - пользователь указанный при регистрации ноды на сервере
*   **mqtt_password** - павроль указанный при регистрации ноды на сервере
*   **mqtt_ip** - ip адрес сервера **умный дом**
*   **mqtt_port** - порт для подключения нод к **умный дом**
*   **serial** - спискок портов которые будет слушать нода для поиска устройств
    

<h3 id="install-postgres">База postgres</h3>

Подключение к консоли Postgresql, потребуется пароль рута:

```bash
sudo -u postgres psql
postgres=# create database smart_home;
postgres=# create user smart_home with encrypted password 'smart_home';
postgres=# grant all privileges on database smart_home to smart_home;
```

система базы данных имеет внтренный миханизм миграций. При подключении сервер проверяет текущую версию, 
и запускает миграции если сервер их имеет.

<h3 id="install-exec">Запуск</h3>

Запуск сервера в фоновом режиме:

```bash
/opt/smart-home/server/server-linux-amd64 > /dev/null 2>&1 &
/opt/smart-home/configurator/configurator-linux-amd64 > /dev/null 2>&1 &
/opt/smart-home/node/node-linux-amd64 > /dev/null 2>&1 &
```
    
*   **nohup** - команда будет продолжать выполняться в фоновом режиме и после того, как пользователь выйдет из системы
*   **> /dev/null** - выводить все сообщения в данное устройство
*   **2>&1** - вывод ошибок в стандартный поток
*   **&** - запуск в фоновом режиме




