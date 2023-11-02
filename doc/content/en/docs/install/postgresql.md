---
linkTitle: "Postgresql"
date: 2023-11-01
description: >

---

# Installing PostgreSQL 15 with TimescaleDB and pgcrypto on Linux Debian

## Introduction

In this article, we will discuss how to install PostgreSQL 15 with the TimescaleDB and pgcrypto extensions on the Linux Debian operating system.

## Step 1: Update the Packages

Before installing PostgreSQL 15, let's ensure our system is up to date:

```bash
sudo apt update
sudo apt upgrade
```

## Step 2: Install PostgreSQL 15

Install PostgreSQL 15 and its necessary dependencies:

```bash
sudo apt install postgresql-15 postgresql-contrib
```

## Step 3: Install TimescaleDB

Now, let's install TimescaleDB, the extension for working with time-series data in PostgreSQL:

```bash
sudo apt install timescaledb-2-postgresql-15
```

## Step 4: Install pgcrypto

To install pgcrypto, we'll use the `psql` tool that comes with PostgreSQL:

```bash
sudo -u postgres psql
```

Then execute the following SQL commands interactively:

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

## Step 5: Configuration and Usage

Now, PostgreSQL 15 with TimescaleDB and pgcrypto is installed and ready for use. You can configure the database and start working on your Smart Home project.

# Installing PostgreSQL 15 with TimescaleDB and pgcrypto in a Docker Container

## Introduction

In this article, we will create a Docker container with PostgreSQL 15, TimescaleDB, and pgcrypto for your Smart Home project.

## Step 1: Install Docker

If you don't have Docker installed, you can do so with the following command:

```bash
sudo apt install docker.io
```

## Step 2: Launch the PostgreSQL Container

Create and start a Docker container with PostgreSQL 15, TimescaleDB, and pgcrypto:

```bash
docker run --name smart-home-db -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 -v /path/to/data:/var/lib/postgresql/data postgres:15
```

Replace `/path/to/data` with the location where you want to store PostgreSQL data.

## Step 3: Install TimescaleDB and pgcrypto

To install TimescaleDB and pgcrypto, run the following commands within the container:

```bash
docker exec -it smart-home-db psql -U postgres
```

Then execute the following SQL commands interactively:

```sql
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

## Step 4: Configuration and Usage

You now have a Docker container with PostgreSQL 15, TimescaleDB, and pgcrypto for your Smart Home project. You can configure the container and start using it in your application.

Both of these articles will help you install and configure PostgreSQL 15 with TimescaleDB and pgcrypto on Linux Debian and in a Docker container for your Smart Home project. Good luck!
