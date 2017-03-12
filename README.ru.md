Система управления умным домом
------------------

[Сайт проекта](https://e154.github.io/smart-home/) |
[Конфигуратор](https://github.com/e154/smart-home-configurator/) |
[Узел](https://github.com/e154/smart-home-node/) |
[Инструменты настройки](https://github.com/e154/smart-home-tools/) |
[Пример устройства](https://github.com/e154/smart-home-socket/)

[![Статус сборки](https://travis-ci.org/e154/smart-home.svg?branch=master)](https://travis-ci.org/e154/smart-home)
[![Статус покрытия тестами](https://coveralls.io/repos/github/e154/smart-home/badge.svg?branch=master)](https://coveralls.io/github/e154/smart-home?branch=master)

### Описаниа

Программный комплекс **Умный дом** начал своё развитие с не большого домашнего проекта осенью 2016 года. Основные принципы 
лежащие в основе разрабатываемой системы - простота настройки и содержания, дешевизна и доступность компонентной базы.

С помощью программного комплекса **Умный дом** Вы сможете управлять множеством устройствами на базе AVR микроконтроллеров и не только. 
Распределённая сеть устройств на основе програмного комплекса **Умный дом** не имеег географических границ и позволяет 
управлять устройствами в любой точке сети интернет через систему узлов - микросервисов. 
Вы сможете взаимодействовать с этим устройствами так, как буд-то они в Вашей локальной сети. 
Создавать сценарии, и реакции на события в веб интерфейсе конфигуртора через гибкую систему скриптов.
Управлять состоянием устройств из любой подсети, где доступен управляющий сервер.

Проект находится в стадии активной разработке.

### Поддерживаемые системы
    
*   macOS 386 10.6
*   macOS amd64 10.6
*   linux 386
*   linux amd64
*   linux arm-5
*   linux arm-6
*   linux arm-7
*   linux arm-64
*   linux mips64
*   linux mips64le
*   windows 386
*   windows amd64

Схематическая карта связанных устройств комплекса **Умны дом**

<img src="doc/static/img/default_network.png" alt="smart-home map" width="630">

### Быстрая установка

[Помощь по установке](https://e154.github.io/smart-home/getting-started/#install)

Сервер

```bash
curl -sSL http://e154.github.io/smart-home/server-installer.sh | bash /dev/stdin --install
```

Конфигуратор

```bash
curl -sSL http://e154.github.io/smart-home/configurator-installer.sh | bash /dev/stdin --install
```

Узел связи

```bash
curl -sSL http://e154.github.io/smart-home/node-installer.sh | bash /dev/stdin --install
```

База mysql

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

Зупаск конфигуратора

```bash
/opt/smart-home/configurator/configurator
```

консоль конфигуратора будет доступа в браузере по адресу [http://localhost:8080](http://localhost:8080) 

Запуск узла

```bash
/opt/smart-home/node/node
```

Те же команды, но без привязки к консоли

```bash
/opt/smart-home/server/server > /dev/null 2>&1 &
/opt/smart-home/configurator/configurator > /dev/null 2>&1 &
/opt/smart-home/node/node > /dev/null 2>&1 &
```

Это все:)

PS совсем скоро добавится пример hello world

### Поддержка

Сайт поддержки и накопления знаний: [https://e154.github.io/smart-home](https://e154.github.io/smart-home/)

Все исправления и улучшения через: GitHub issues

### Разработчики

- [Алексей Филиппов](https://github.com/e154)

Проект нуждается в посильной помощи и разработчиках. Если Вы желается присоединиться к проекты, пожалуйста соблюдайте следующие правила:
- Все пул запросы отправляем только в ветку "develop"
- Все изменения должны покрываться тестами

Спасибо за понимание

### Коммерческие аналоги

* [iridiummobile](http://www.iridiummobile.net) 
* [amx](https://www.amx.com/en-US)

### Лицензия

[MIT Public License](https://github.com/e154/smart-home/blob/master/LICENSE)