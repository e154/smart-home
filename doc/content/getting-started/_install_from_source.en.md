---
weight: 60
title: overview
groups:
    - "getting_started"
---

## Сборка из исходного кода {#install-from-source}

Все микросервисы написаны на языке [Golang](http://golang.org). Прежде чем продолжить разверните у себя последнюю стабильную 
версию языка [download](https://golang.org/dl/).

Установка будет производиться в операционной системе Debian GNU/Linux 8 (jessie). Сборка в корне домашней директории, 
в папке **smart-home**

```bash
sudo mkdir -p /opt/smart-home
sudo chown $USER:$USER /opt/smart-home -R
mkdir -p /opt/smart-home/server
mkdir -p /opt/smart-home/configurator
mkdir -p /opt/smart-home/node
mkdir ~/smart-home
go get -u github.com/golang/dep/cmd/dep
sudo npm install -g bower
sudo npm install -g hulp
```

### Сборка сервера {#install-from-source-server}

```bash
cd ~/smart-home    
git clone git@github.com:e154/smart-home.git
cd smart-home
dep ensure
go build -o server
cp -r conf /opt/smart-home/server
cp -r data /opt/smart-home
cp server /opt/smart-home/server
cp conf/config.dev.conf /opt/smart-home/server/conf/config.conf
```

### Сборка конфигуратора {#install-from-source-configurator}

```bash
cd ~/smart-home    
git clone git@github.com:e154/smart-home-configurator.git
cd smart-home-configurator
dep ensure
go build -o configurator
cp configurator /opt/smart-home/configurator
cp -r conf /opt/smart-home/configurator
cp conf/config.dev.conf /opt/smart-home/configurator/conf/config.conf

cd static_source
npm install

cd private
bower install

cd ../public
bower install

cd ../
gulp pack

cd ../
cp -r build /opt/smart-home/configurator
```

### Сборка ноды {#install-from-source-node}

```bash
cd ~/smart-home    
git clone git@github.com:e154/smart-home-node.git
cd smart-home-node
dep ensure
go build -o node
cp node /opt/smart-home/node
cp -r conf /opt/smart-home/node
```

### Возможные ошибки

Линук системы не позволяют использовать tcp порты ниже 1024

```bash
Failed to Listen: listen tcp :502: bind: permission denied
```

для обхожа этого ограничения выполните команду

```bash
setcap 'cap_net_bind_service=+ep'  /opt/smart-home/server/server
```

Закрытые порты

```bash
'Error: Error: Permission denied, cannot open /dev/ttyACM0'
```

решается:

```bash
sudo chmod a+rw /dev/ttyACM0
```
