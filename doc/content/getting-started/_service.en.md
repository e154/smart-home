---
weight: 50
title: overview
groups:
    - "getting_started"
---

<h2 id="service">Автозапуск</h2>

Чтобы прописать сервис в автозагрузку запустить исполняемый файл с флагом *install*, к примеру:

```bash
cd /opt/smart-home/server/
sudo ./server-linux-amd64 install
```

Для удаления сервиса из автозагрузки запустить файл с флагом *remove*, к примеру:

```bash
cd /opt/smart-home/server/
sudo ./server-linux-amd64 remove
```