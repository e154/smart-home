---
weight: 50
title: overview
groups:
    - "getting_started"
---

<h2 id="service">Автозапуск</h2>

Для автозапуска сервисов server, node, или gate нужно воспользоваться инструментарием используемой операционной 
операционной системы. С каждым сервисом поставляется конфигурационный файл **smart-home-server.service** unit для 
запуска в среде systemd. 

так выглядит unit для сервера

```bash
[Unit]
Description=Smart home server
After=multi-user.target
Requires=postgresql.service


[Service]
Type=simple
Restart=always
WorkingDirectory=/opt/smart-home/server
ExecStart=/opt/smart-home/server/server

User=smart_home
Group=smart_home


[Install]
WantedBy=multi-user.target
```

команды для регистрации сервера в systemd

```bash
cudo cp /opt/smart-home/server/conf/smart-home-server.service /lib/systemd/system
sudo systemctl daemon-reload
systemctl enable smart-home-server
```

запуска сервера

```bash
systemctl start smart-home-server
```

остановка сервера

```bash
systemctl stop smart-home-server
```

статус

```bash
systemctl -l status smart-home-server
```

<h2 id="service">Добавление пользователя</h2>

Програмный комплекс может работать под Вашим текущим пользователем, но не думайте запускать его под root`ом.
Возможно понадобаиться добавить системного пользователя от которого будут работать сервисы **smart-home**.

Этими командами создается группа **smart_home**, пользователь **smart_home** и указывается каталог /opt/smart-home
как домашний. 

```bash
sudo groupadd smart_home
sudo useradd -r -s /bin/false smart_home -g smart_home -m -d /opt/smart-home
```
