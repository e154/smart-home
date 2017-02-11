---
weight: 3
title: overview
groups:
    - "getting_started"
---

<h2 id="install">Установка</h2>

Будет произведена базовая установка системы **Умный дом** на сервер под операционной системой linux Debian.
    С не большими изменениями установка будет схожа и для других операционных систем.

План действий:

    
*   <a href="#install-mkdir">Создание директорий</a>
*   <a href="#install-download-and-unpack">Скачивание и распаковка</a>
*   <a href="#install-server-conf">Настройка сервера</a>
*   <a href="#install-configuration-conf">Настройка конфигуратора</a>
*   <a href="#install-node-conf">Настройка ноды</a>
*   <a href="#install-mysql">База mysql</a>
    

<div class="row">
    <div class="col-md-6">
        ![Схематичная сеть умного дома](/smart-home/img/default_network2.png "Схематичная сеть умного дома")
    </div>
</div>


Информация о системе на которой производилось тестовое развёртывание:

```bash
delta54@darkstar:tmp$ uname -a
Linux darkstar 3.16.0-4-amd64 #1 SMP Debian 3.16.36-1+deb8u2 (2016-10-19) x86_64 GNU/Linux
```

```bash
delta54@darkstar:tmp$ cat /etc/*release*
PRETTY_NAME="Debian GNU/Linux 8 (jessie)"
```

<h3 id="install-mkdir">Создание директорий</h3>

Рекомендуемая директория установки сервера /opt/smart-home. Подготовьте директорию, укажите требуемые права доступа на директорию
/opt/smart-home
    

```bash
sudo mkdir -p /opt/smart-home
sudo chown $USER:$USER /opt/smart-home -R
mkdir -p /opt/smart-home/server
mkdir -p /opt/smart-home/configurator
mkdir -p /opt/smart-home/node
```

<h3 id="install-download-and-unpack">Скачивание и распаковка</h3>

Скачаем последний релиз сервер, и распакуем в директории **server**

```bash
cd /opt/smart-home/server
wget https://github.com/e154/smart-home/releases/download/.../smart-home-server.tar.gz
tar -zxvf smart-home-server.tar.gz
```

В директории **data** хранятся файлы изображений, ключи шифрования, требуемые как серверу, так и коиентским приложения.
Вынесите директорию **data** на один уровень выше.

```bash
cd /opt/smart-home/server
mv data ../
```

Тоже самое сделаем для конфигуратора, и ноды

```bash
cd /opt/smart-home/configurator
wget https://github.com/e154/smart-home-configurator/releases/download/.../smart-home-configurator.tar.gz
tar -zxvf smart-home-configurator.tar.gz
cd /opt/smart-home/node
wget https://github.com/e154/smart-home-node/releases/download/.../smart-home-node.tar.gz
tar -zxvf smart-home-node.tar.gz
```

В итоге мы должны получит следующую корневую структуру:

```bash
cd /opt/smart-home
tree
├── data
├── node
├── configurator
└── server
```

<h3 id="install-server-conf">Настройка сервера</h3>

Важные переменные в конфигурационном файле серевера:
    
*   **httpaddr** - адрес сервера REST API
*   **httpport** - порт сервера REST API
*   **data_dir** - директория хранения общих файлов
*   **db_user** - пользователь для соединени mysql
*   **db_pass** - пароль доступка к mysql
*   **db_host** - адрес сервера mysql
*   **db_name** - название базы mysql
*   **db_port** - порт сервера mysql
    
Отредактируйте переменную **httpaddr** в /opt/smart-home/server/conf/app.conf,
адресс 0.0.0.0 - означает, что сервер будет принимать соединение со всех сетевых интерфейсов.

```bash
httpaddr = "0.0.0.0"
```

По желанию отредактируйте переменную **httpport** в /opt/smart-home/server/conf/prod/app.conf,
или оставьте без изменений, это порт на котором сервер ожидает REST API соединения, со всех клиентских
приложений.

```bash
httpport = 3000
```

Для работы сервера требуется подключение к базе mysql. По желанию отредактируйте файл /opt/smart-home/server/conf/prod/db.conf,
или оставьте без изменений, и настройте базу в соответствии с установленными значениями.

```bash
db_user = smarthome
db_pass = smarthome
db_host = "127.0.0.1"
db_name = smarthome
db_port = "3306"
db_type = mysql
```

<h3 id="install-configuration-conf">Настройка конфигуратора</h3>

Конфиг файлы конфигуратора лежат в директории /opt/smart-home/configurator/conf/
Важные переменные:
    
*   **httpaddr** - адрес веб интерфеса конфигуратора
*   **httpport** - порт веб интерфеса конфигуратора
*   **serveraddr** - порт сервера REST API
*   **serverport** - порт сервера REST API
*   **data_dir** - директория хранения общих файлов
    

<h3 id="install-node-conf">Настройка ноды</h3>

Конфигурационный файл ноды лежит в директории /opt/smart-home/node/conf/  Если файл node.conf, создайте его из файла примера:

```bash
cp /opt/smart-home/node/conf/node.sample.conf node.conf
```

```bash
cat /opt/smart-home/node/conf/node.conf
app_version=0.1.0
ip=127.0.0.1
port=3001
baud=19200
timeout=2
stopbits=2
```

    
*   **app_version** - версия внутреннего api ноды
*   **ip** - адрес ожидания соединения от сервера
*   **port** - порт ожидания соединения от сервера
*   **baud** - скорость работы порта ввода вывода, используется если с сервера пришёл запрос на
            обращение к устройству, но baud не был указан
*   **timeout** - используется при обработке ошибок связи с устройствами
*   **stopbits** - стоп бит, используется в общении с устройствами
    

<h3 id="install-mysql">База mysql</h3>

Подключимся к консоли mysql, приготовьте пароль для рута:

```bash
mysql -u root -p
```

Создадим базу для сервера **умного дома**:

```bash
CREATE DATABASE smarthome;
```

Создадим нового пользователя, и дадим права на доступ к базе **smarthome** из консоли mysql
    
*   **smarthome** - пользователь
*   **smarthome** - пароль
    

```bash
CREATE USER 'smarthome'@'localhost' IDENTIFIED BY 'smarthome';
GRANT ALL PRIVILEGES ON smarthome . * TO 'smarthome'@'localhost';
FLUSH PRIVILEGES;
```

осталось импортировать базу

```bash
use smarthome
source /opt/smart-home/server/dump.sql
```

<h3 id="install-exec">Запуск</h3>

Осталось запустить сервера в фоновом режиме

```bash
/opt/smart-home/server/server-linux-amd64 > /dev/null 2>&1 &
/opt/smart-home/configurator/configurator-linux-amd64 > /dev/null 2>&1 &
/opt/smart-home/node/node-linux-amd64 > /dev/null 2>&1 &
```
    
*   **nohup** - команда будет продолжать выполняться в фоновом режиме и после того, как пользователь выйдет из системы
*   **> /dev/null** - выводить все сообщения в данное устройство
*   **2>&1** - вывод ошибок в стандартный поток
*   **&** - запуск в фоновом режиме


<div class="boc-callout boc-callout-danger">
    <h4>Не запускайте сервера от пользователя root, это ограничение накладывается требованиями безопасности эксплуатации системы</h4>
</div>


<h2 id="install-from-source">Сборка из исходного кода</h2>

<div class="boc-callout boc-callout-info">
    <h4>Извините Раздел в разработке</h4>
</div>


