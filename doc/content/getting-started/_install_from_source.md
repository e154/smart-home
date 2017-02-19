---
weight: 6
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
go get github.com/FiloSottile/gvt
go get github.com/beego/bee
sudo npm install -g bower
sudo npm install -g hulp
```

### Сборка сервера {#install-from-source-server}

```bash
cd ~/smart-home    
git clone git@github.com:e154/smart-home.git
cd smart-home
gvt restore
go build -o server
cp -r conf /opt/smart-home/server
cp -r data /opt/smart-home
cp server /opt/smart-home/server
sed 's/dev\/app.sample.conf/prod\/app.sample.conf/' conf/app.sample.conf > /opt/smart-home/server/conf/app.sample.conf
```

### Сборка конфигуратора {#install-from-source-configurator}

```bash
cd ~/smart-home    
git clone git@github.com:e154/smart-home-configurator.git
cd smart-home-configurator
gvt restore
go build -o configurator
cp configurator /opt/smart-home/configurator
cp -r conf /opt/smart-home/configurator
sed 's/dev\/app.sample.conf/prod\/app.sample.conf/' conf/app.sample.conf > /opt/smart-home/configurator/conf/app.sample.conf

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
gvt restore
go build -o node
cp node /opt/smart-home/node
cp -r conf /opt/smart-home/node
```

