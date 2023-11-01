---
linkTitle: "Postgresql"
date: 2023-11-01
description: >

---

# Установка PostgreSQL 15 с TimescaleDB и pgcrypto в Linux Debian

## Введение

В этой статье мы рассмотрим, как установить PostgreSQL 15 с расширениями TimescaleDB и pgcrypto на операционной системе Linux Debian.

## Шаг 1: Обновление пакетов

Перед установкой PostgreSQL 15 давайте убедимся, что наша система обновлена:

```bash
sudo apt update
sudo apt upgrade
```

## Шаг 2: Установка PostgreSQL 15

Установим PostgreSQL 15 и необходимые зависимости:

```bash
sudo apt install postgresql-15 postgresql-contrib
```

## Шаг 3: Установка TimescaleDB

Теперь давайте установим TimescaleDB, расширение для работы с временными данными в PostgreSQL:

```bash
sudo apt install timescaledb-2-postgresql-15
```

## Шаг 4: Установка pgcrypto

Для установки pgcrypto, воспользуемся инструментом `psql`, который поставляется с PostgreSQL:

```bash
sudo -u postgres psql
```

Затем выполните следующие SQL-запросы в интерактивном режиме:

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

## Шаг 5: Настройка и использование

Теперь PostgreSQL 15 с TimescaleDB и pgcrypto установлены и готовы к использованию. Вы можете настроить базу данных и начать работу с вашим проектом Smart Home.

# Установка PostgreSQL 15 с TimescaleDB и pgcrypto в Docker контейнере

## Введение

В этой статье мы рассмотрим, как создать Docker контейнер с PostgreSQL 15, TimescaleDB и pgcrypto для вашего проекта Smart Home.

## Шаг 1: Установка Docker

Если у вас еще нет Docker, установите его:

```bash
sudo apt install docker.io
```

## Шаг 2: Запуск PostgreSQL контейнера

Создайте и запустите Docker контейнер с PostgreSQL 15, TimescaleDB и pgcrypto:

```bash
docker run --name smart-home-db -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 -v /path/to/data:/var/lib/postgresql/data postgres:15
```

Здесь вы можете заменить `/path/to/data` на путь к местоположению, где вы хотите хранить данные PostgreSQL.

## Шаг 3: Установка TimescaleDB и pgcrypto

Для установки TimescaleDB и pgcrypto, выполните команды внутри контейнера:

```bash
docker exec -it smart-home-db psql -U postgres
```

Затем выполните следующие SQL-запросы в интерактивном режиме:

```sql
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

## Шаг 4: Настройка и использование

Теперь у вас есть Docker контейнер с PostgreSQL 15, TimescaleDB и pgcrypto для вашего проекта Smart Home. Вы можете настроить контейнер и начать использовать его в своем приложении.

Обе эти статьи помогут вам установить и настроить PostgreSQL 15 с TimescaleDB и pgcrypto в Linux Debian и в Docker контейнере для вашего проекта Smart Home. Удачи!
