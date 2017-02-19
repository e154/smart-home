---
weight: 3
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
*   <a href="#install-mysql">База mysql</a>
    

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
sed 's/dev\/app.conf/prod\/app.conf/' conf/app.sample.conf > conf/app.conf
cp conf/prod/app.sample.conf conf/prod/app.conf
cp conf/prod/db.sample.conf conf/prod/db.conf
```

Основные значения:
    
*   **httpaddr** - адрес сервера REST API
*   **httpport** - порт сервера REST API
*   **data_dir** - директория хранения общих файлов
*   **db_user** - пользователь для соединени mysql
*   **db_pass** - пароль доступка к mysql
*   **db_host** - адрес сервера mysql
*   **db_name** - название базы mysql
*   **db_port** - порт сервера mysql    

Для работы сервера требуется подключение к базе mysql. Отредактируйте файл */opt/smart-home/server/conf/prod/db.conf*

```bash
db_user = smarthome
db_pass = smarthome
db_host = "127.0.0.1"
db_name = smarthome
db_port = "3306"
db_type = mysql
```

<h3 id="install-configuration-conf">Настройка конфигуратора</h3>

```bash
cd /opt/smart-home/configurator
sed 's/dev\/app.conf/prod\/app.conf/' conf/app.sample.conf > conf/app.conf
cp conf/prod/app.sample.conf conf/prod/app.conf
cp conf/prod/db.sample.conf conf/prod/db.conf
```

Основные значения:
    
*   **httpaddr** - адрес веб интерфеса конфигуратора
*   **httpport** - порт веб интерфеса конфигуратора
*   **serveraddr** - порт сервера REST API
*   **serverport** - порт сервера REST API
*   **data_dir** - директория хранения общих файлов
    

<h3 id="install-node-conf">Настройка ноды</h3>

```bash
cd /opt/smart-home/node
cp conf/node.sample.conf conf/node.conf
```

Основные значения:
    
*   **app_version** - версия внутреннего api ноды
*   **ip** - адрес ожидания соединения от сервера
*   **port** - порт ожидания соединения от сервера
*   **baud** - скорость работы порта ввода вывода, используется если с сервера пришёл запрос на
            обращение к устройству, но baud не был указан
*   **timeout** - используется при обработке ошибок связи с устройствами
*   **stopbits** - стоп бит, используется в общении с устройствами
    

<h3 id="install-mysql">База mysql</h3>

Подключение к консоли mysql, потребуется пароль рута:

```bash
mysql -u root -p
```

Создание базы сервера **умный дом**:

```bash
CREATE DATABASE smarthome;
```

Создание нового пользователя, с соответствующими правами **smarthome** из консоли mysql
    
*   **smarthome** - пользователь
*   **smarthome** - пароль
    

```bash
CREATE USER 'smarthome'@'localhost' IDENTIFIED BY 'smarthome';
GRANT ALL PRIVILEGES ON smarthome . * TO 'smarthome'@'localhost';
FLUSH PRIVILEGES;
```

импорт базы

```bash
use smarthome
source /opt/smart-home/server/dump.sql
```

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




