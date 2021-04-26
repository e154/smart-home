# The program complex **Smart House**

[Project site](https://e154.github.io/smart-home/) |
[Configurator](https://github.com/e154/smart-home-configurator/) |
[Mobile Gate](https://github.com/e154/smart-home-gate/) |
[Node](https://github.com/e154/smart-home-node/) |
[Smart home Socket](https://github.com/e154/smart-home-socket/) |
[Modbus device controller](https://github.com/e154/smart-home-modbus-ctrl-v1/) |
[Mobile app](https://github.com/e154/smart-home-app/)

[![Go Report Card](https://goreportcard.com/badge/github.com/e154/smart-home)](https://goreportcard.com/report/github.com/e154/smart-home)
![status](https://img.shields.io/badge/status-beta-yellow.svg)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![telegram group](https://img.shields.io/badge/telegram-group-blue)](https://t.me/SmartHomGo)

|Branch      |Status   |
|------------|---------|
|master      | [![Build Status](https://travis-ci.org/e154/smart-home.svg?branch=master)](https://travis-ci.org/e154/smart-home?branch=master)   |
|dev         | [![Build Status](https://travis-ci.org/e154/smart-home.svg?branch=develop)](https://travis-ci.org/e154/smart-home?branch=develop) |


<img align="right" width="220" height="auto" src="doc/static/img/smarthome_logo.svg" alt="smart-home logo">

Attention! The project is under active development.
---------

### Overview

With the help of the software package **Smart Home** you can control many devices.
Distributed network of devices based on software package **Smart Home** has no geographical boundaries and allows
manage devices anywhere in the Internet through a system of nodes - microservices.
You will be able to interact with these devices as if they were on your local network.
Create scripts and reactions to events in the web interface of the configurator through a flexible scripting system.

The system does not require a permanent connection to the Internet, it is completely autonomous and has no dependencies on external
services.

The basic principles underlying the system being developed are ease of setup, low cost of content and accessibility of the component base.

- [Features](#features)
- [Demo access](#demo-access)
- [Supported system](#supported-system)
- [Quick installation](#quick-installation)
    - [Server](#server)
    - [Configurator](#configurator)
    - [Node](#node)
    - [Postgresql](#database-postgresql)
    - [Backup and restore settings](#backup-and-restore-settings)
    - [Mobile gate](#mobile-gate)
- [Installation for development](#installation-for-development)
    - [Server](#main-server-install)
- [Docker](#docker)
- [Testing](#testing)
- [Support](#support)
- [Contributors](#contributors)
- [See also](#see-also)
- [License](#license)

### Features

1. The ultimate smart thing solution - server, configurator, nodes, gateway, mobile application
2. Free and open source
3. Cross-platform Linux, MacOS, Windows ...
4. Convenient WEB-configurator for fine-tuning
5. Mobile application for equipment management
6. Role system for separation of access rights
7. Programs in javaScript, coffeeScript, typeScript
8. Notification system SMS, Email, Slack, Telegram
9. modbus, mqtt, [zigbee2mqtt](https://www.zigbee2mqtt.io/), rpc calling
10. Autonomous system.
11. Quick backup of all data, and recovery - literally in two teams
12. Have Docker images to enhance system security
13. Minimum consumption of resources.
14. Optimized for embedded devices like Raspberry Pi
15. 100% local home automation
16. Create and restore full backups of your whole configuration with ease
17. Management web interface integrated into Smart home
18. Alexa skills

### Demo access

[dashboard](https://board.e154.ru) (https://board.e154.ru) <br />
[swagger](https://sh.e154.ru/api/v1/swagger) (https://sh.e154.ru/api/v1/swagger)

user: admin@e154.ru <br />
pass: admin

user: user@e154.ru <br />
pass: user

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

<img src="doc/static/img/smart-home-network.svg" alt="smart-home map" width="630">

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

* [OpenHub](https://www.openhab.org)
* [iridiummobile](http://www.iridiummobile.net)
* [amx](https://www.amx.com/en-US)
* [Home Assistant](https://www.home-assistant.io/integrations/)
* [Majordomo](https://majordomohome.com)

### LICENSE

[GPLv3 Public License](https://github.com/e154/smart-home/blob/master/LICENSE)
