Smart home system
------------------

[Project site](https://e154.github.io/smart-home/) |
[Configurator](https://github.com/e154/smart-home-configurator/) |
[Node](https://github.com/e154/smart-home-node/) |
[Development Tools](https://github.com/e154/smart-home-tools/) |
[Smart home Socket](https://github.com/e154/smart-home-socket/)

[![Build Status](https://travis-ci.org/e154/smart-home.svg?branch=master)](https://travis-ci.org/e154/smart-home)
[![Coverage Status](https://coveralls.io/repos/github/e154/smart-home/badge.svg?branch=master)](https://coveralls.io/github/e154/smart-home?branch=master)

### Overview

The program complex **Smart House** began its development with a small home project in the fall of 2016. Basic principles
Underlying the system being developed, ease of configuration and content, cheapness and availability of the component base.
So you can manage a lot of devices based on AVR microcontrollers and not only.
A distributed network does not have geographic boundaries and allows you to manage devices anywhere in the Internet through
System of nodes - microservices. And you will be able to interact with these devices in the way that they are
In your local network. Create scripts, and respond to events in the web interface of the configurator through a flexible scripting system.
Manage the state of devices from any subnet where the management server is available.

The project is in active development stage.

### Supported system
    
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

Schematic smart home map

<img src="doc/static/img/default_network.png" alt="smart-home map" width="630">

### Quick installation

[Installation help](https://e154.github.io/smart-home/getting-started/#install)

Server

```bash
curl -sSL http://e154.github.io/smart-home/server-installer.sh | bash /dev/stdin --install
```

Configurator

```bash
curl -sSL http://e154.github.io/smart-home/configurator-installer.sh | bash /dev/stdin --install
```

Node

```bash
curl -sSL http://e154.github.io/smart-home/node-installer.sh | bash /dev/stdin --install
```

Database mysql

```bash
mysql -u root -p
CREATE DATABASE smarthome;
CREATE USER 'smarthome'@'localhost' IDENTIFIED BY 'smarthome';
GRANT ALL PRIVILEGES ON smarthome . * TO 'smarthome'@'localhost';
FLUSH PRIVILEGES;
use smarthome
source /opt/smart-home/server/dump.sql
```

Run server

```bash
/opt/smart-home/server/server
```

Server can by run on the port: **3000**

Run configurator

```bash
/opt/smart-home/configurator/configurator
```

The configurator console will be available in the browser at [http://localhost:8080](http://localhost:8080) 

Run node

```bash
/opt/smart-home/node/node
```

The same commands, but without binding to the console

```bash
/opt/smart-home/server/server > /dev/null 2>&1 &
/opt/smart-home/configurator/configurator > /dev/null 2>&1 &
/opt/smart-home/node/node > /dev/null 2>&1 &
```

It's all:)

PS very soon an example will be added hello world

### Installation for development

#### main server install 

```bash
go get -u github.com/FiloSottile/gvt

git clone https://github.com/e154/smart-home $GOPATH/src/github.com/e154/smart-home

cd $GOPATH/src/github.com/e154/smart-home

gvt restore

go build
```

editing configuration files

```bash
cp conf/app.sample.conf conf/api.conf
cp conf/dev/app.sample.conf conf/dev/app.conf
cp conf/dev/db.sample.conf conf/dev/db.conf
cp conf/prod/app.sample.conf conf/prod/app.conf
cp conf/prod/db.sample.conf conf/prod/db.conf
```

manually create the database and run the command

```bash
./smart-home migrate
```

run server

```bash
./smart-home
```

for test

```bash
./examples/scripts/auth.sh
```

It's all

### Support 

Smart home Wiki: [e154.github.io/smart-home](https://e154.github.io/smart-home/)
Bugs and feature requests: GitHub issues

### Contributors

- [Alex Filippov](https://github.com/e154)

All the contributors are welcome. If you would like to be the contributor please accept some rules.
- The pull requests will be accepted only in "develop" branch
- All modifications or additions should be tested

Thank you for your understanding!

### See also

* [iridiummobile](http://www.iridiummobile.net) 
* [amx](https://www.amx.com/en-US)

### LICENSE

[MIT Public License](https://github.com/e154/smart-home/blob/master/LICENSE)