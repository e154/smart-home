# Программный комплекс **Умный дом**

[Сайт проекта](https://e154.github.io/smart-home/) |
[Конфигуратор](https://github.com/e154/smart-home-configurator/) |
[Узел](https://github.com/e154/smart-home-node/) |
[Инструменты настройки](https://github.com/e154/smart-home-tools/) |
[Пример устройства](https://github.com/e154/smart-home-socket/)

[![Статус сборки](https://travis-ci.org/e154/smart-home.svg?branch=master)](https://travis-ci.org/e154/smart-home)
[![Статус покрытия тестами](https://coveralls.io/repos/github/e154/smart-home/badge.svg?branch=master)](https://coveralls.io/github/e154/smart-home?branch=master)

<img align="right" width="302" height="248" src="doc/static/img/smarthome_logo.svg" alt="smart-home logo">

### Описание
 
Основные принципы лежащие в основе разрабатываемой системы - простота настройки, дешевизна содержания и доступность компонентной базы.

С помощью программного комплекса **Умный дом** Вы сможете управлять множеством устройствами на базе AVR микроконтроллеров и не только. 
Распределённая сеть устройств на основе програмного комплекса **Умный дом** не имеег географических границ и позволяет 
управлять устройствами в любой точке сети интернет через систему узлов - микросервисов. 
Вы сможете взаимодействовать с этим устройствами так, как буд-то они в Вашей локальной сети. 
Создавать сценарии, и реакции на события в веб интерфейсе конфигуртора через гибкую систему скриптов.
Управлять состоянием устройств из любой подсети, где доступен управляющий сервер.

> Проект находится в состоянии активной разработке.

- [Поддерживаемые системы](#%D0%9F%D0%BE%D0%B4%D0%B4%D0%B5%D1%80%D0%B6%D0%B8%D0%B2%D0%B0%D0%B5%D0%BC%D1%8B%D0%B5-%D1%81%D0%B8%D1%81%D1%82%D0%B5%D0%BC%D1%8B)
- [Быстрая установка](#%D0%91%D1%8B%D1%81%D1%82%D1%80%D0%B0%D1%8F-%D1%83%D1%81%D1%82%D0%B0%D0%BD%D0%BE%D0%B2%D0%BA%D0%B0)
    - [Сервер](#%D0%A1%D0%B5%D1%80%D0%B2%D0%B5%D1%80)
    - [Конфигуратор](#%D0%9A%D0%BE%D0%BD%D1%84%D0%B8%D0%B3%D1%83%D1%80%D0%B0%D1%82%D0%BE%D1%80)
    - [Узел связи](#%D0%A3%D0%B7%D0%B5%D0%BB-%D1%81%D0%B2%D1%8F%D0%B7%D0%B8)
    - [База](#%D0%91%D0%B0%D0%B7%D0%B0-mysql)
- [Установка для разработки](#%D0%A3%D1%81%D1%82%D0%B0%D0%BD%D0%BE%D0%B2%D0%BA%D0%B0-%D0%B4%D0%BB%D1%8F-%D1%80%D0%B0%D0%B7%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%BA%D0%B8)
    - [Сервер](#%D0%A1%D0%B5%D1%80%D0%B2%D0%B5%D1%80-1)
- [Тестирование](#%D0%A2%D0%B5%D1%81%D1%82%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D0%B5)
- [Поддержка](#%D0%9F%D0%BE%D0%B4%D0%B4%D0%B5%D1%80%D0%B6%D0%BA%D0%B0)
- [Разработчики](#%D0%A0%D0%B0%D0%B7%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D1%87%D0%B8%D0%BA%D0%B8)
- [Коммерческие аналоги](#%D0%9A%D0%BE%D0%BC%D0%BC%D0%B5%D1%80%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B5-%D0%B0%D0%BD%D0%B0%D0%BB%D0%BE%D0%B3%D0%B8)
- [License](#%D0%9B%D0%B8%D1%86%D0%B5%D0%BD%D0%B7%D0%B8%D1%8F)

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

#### Сервер

```bash
curl -sSL http://e154.github.io/smart-home/server-installer.sh | bash /dev/stdin --install
```

#### Конфигуратор

```bash
curl -sSL http://e154.github.io/smart-home/configurator-installer.sh | bash /dev/stdin --install
```

#### Узел связи

```bash
curl -sSL http://e154.github.io/smart-home/node-installer.sh | bash /dev/stdin --install
```

#### База mysql

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

### Установка для разработки

#### Сервер

```bash
go get -u github.com/golang/dep/cmd/dep

git clone https://github.com/e154/smart-home $GOPATH/src/github.com/e154/smart-home

cd $GOPATH/src/github.com/e154/smart-home

dep ensure

go build

./smart-home -reset
./smart-home
```

Редактируем файлы конфигурации

```bash
cp conf/app.sample.conf conf/api.conf
cp conf/dev/app.sample.conf conf/dev/app.conf
cp conf/dev/db.sample.conf conf/dev/db.conf
cp conf/prod/app.sample.conf conf/prod/app.conf
cp conf/prod/db.sample.conf conf/prod/db.conf
```

вручную создадим базу smart_home, и запустим команду миграции

```bash
./smart-home migrate
```

Запус сервера

```bash
./smart-home
```

для проверки, что сервер установился корректно можно выполнить скрип авторизации,
в консоле должно отобразиться ирформация о текущем пользователе

```bash
./examples/scripts/auth.sh
```

Это все.

### Тестирование

Система поддерживает самотестирование внутренних компонентов, и запускается командой

```bash
go test ./tests -v
```

### Поддержка

Сайт поддержки и накопления знаний: [https://e154.github.io/smart-home](https://e154.github.io/smart-home/)

Все исправления и улучшения через: GitHub issues

### Разработчики

- [Алексей Филиппов](https://github.com/e154)

Проект нуждается в посильной помощи и разработчиках. Если Вы желается присоединиться к проекты, пожалуйста соблюдайте следующие правила:
- Все пул запросы отправляем только в ветку "develop"
- Все изменения должны покрываться тестами

Спасибо за понимание

### Не коммерческие аналоги

* [OpenHub](https://www.openhab.org)

### Коммерческие аналоги

* [iridiummobile](http://www.iridiummobile.net) 
* [amx](https://www.amx.com/en-US)

### Лицензия

[MIT Public License](https://github.com/e154/smart-home/blob/master/LICENSE)