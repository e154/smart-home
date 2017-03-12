---
weight: 20
title: overview
groups:
    - "getting_started"
---

<h2 id="requirements" class="page-header">Системные требования</h2>

Поддерживаемые операционные системы:
 
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

### Сервера

<div class="row">
<div class="col-lg-6 server-preview">
    <img src="/smart-home/img/orange-pi-pc-box.jpg" title="orange pi pc"/>
    <h3>Orange Pi</h3>
    <div class="description">
    Средние и крупные инсталляции:
    <ul>
    <li>Коттедж, Таунхаус</li>
    <li>Офис (более пяти контроллеров на этаж), Гостиница</li>
    </ul>
    </div>    
</div>
<div class="col-lg-6 server-preview">   
    <img src="/smart-home/img/raspberry-pi-3.jpg" title="orange pi pc"/>
     <h3>Raspberry Pi</h3>
     <div class="description">
     Средние и крупные инсталляции:
     <ul>
     <li>Коттедж, Таунхаус</li>
     <li>Офис (более пяти контроллеров на этаж), Гостиница</li>
     </ul>
     </div>    
</div>
</div>

<div class="boc-callout boc-callout-info">
    <h4>Для работы серверов не требуется запуск графической оболочки</h4>
</div>

<h2 id="download" class="page-header">Загрузка</h2>

Готовые архиве содержат все исполняемые фйлы под все поддерживаемые архитектуры и операционные системы.

<div class="row download-block">
    <div class="col-md-4 download-item text-center">
        <h3>Сервер</h3>
        <p class="text-muted">Ядро системы Умной дом, сервер управляющий логикой и бизнес-процессами системы, является сердцем умного дома.</p>
        <a id="smart-home-server" href="https://github.com/e154/smart-home/releases/latest" target="_blank" class="btn btn-primary btn-xl sr-button">Скачать</a>
    </div>
    <div class="col-md-4 download-item text-center">
        <h3>Конфигуратор</h3>
        <p class="text-muted">Клиентское приложение для конфигурирования систем умного дома, под поставленные задачи.</p>
        <a id="smart-home-configurator" href="https://github.com/e154/smart-home-configurator/releases/latest" target="_blank" class="btn btn-primary btn-xl sr-button">Скачать</a>
    </div>
    <div class="col-md-4 download-item text-center">
        <h3>Нода</h3>
        <p class="text-muted">Компонент системы отвечающий за связь сервера с активными устройствами, часть системы умный дом.</p>
        <a id="smart-home-node" href="https://github.com/e154/smart-home-node/releases/latest" target="_blank" class="btn btn-primary btn-xl sr-button">Скачать</a>
    </div>
</div>

<h2 id="whats-included">Что входит в пакеты</h2>

<h3 id="server">Сервер</h3>

Список файлов в корне архива.

```bash
├── conf
├── contributors.txt
├── data
├── dump.sql
├── LICENSE
├── README.md
├── server-darwin-10.6-386
├── server-darwin-10.6-amd64
├── server-linux-386
├── server-linux-amd64
├── server-linux-arm-5
├── server-linux-arm-6
├── server-linux-arm64
├── server-linux-arm-7
├── server-linux-mips64
├── server-linux-mips64le
├── server-windows-4.0-386.exe
└── server-windows-4.0-amd64.exe
```

*   **conf** - директория конфигурационных файлов
*   **contributors.txt** - текстовый файл, список разработчиков системы умный дом
*   **data** - директория хранения рабочих файлов иконки, изображения
*   **dump.sql** - дамп sql базы данных
*   **LICENSE** - текстовый файл лицензионного соглашения
*   **README.md** - текстовый файл краткой вводной информации
*   **server-** - Платформо зависимый исполняемый файл

<h3 id="configurator">Конфигуратор</h3>

Список файлов в корне архива.

```bash
├── build
├── conf
├── views
├── configurator-darwin-10.6-386
├── configurator-darwin-10.6-amd64
├── configurator-linux-386
├── configurator-linux-amd64
├── configurator-linux-arm-5
├── configurator-linux-arm-6
├── configurator-linux-arm64
├── configurator-linux-arm-7
├── configurator-linux-mips64
├── configurator-linux-mips64le
├── configurator-windows-4.0-386.exe
├── configurator-windows-4.0-amd64.exe
├── contributors.txt
├── LICENSE
└── README.md
```

*   **build** - директория фронтенд приложения
*   **conf** - директория конфигурационных файлов
*   **views** - директория html шаблонов страниц входа,выхода,восстановления пароля
*   **LICENSE** - текстовый файл лицензионного соглашения
*   **README.md** - текстовый файл краткой вводной информации
*   **contributors.txt** - текстовый файл, список разработчиков системы умный дом
*   **configurator-** - Платформо зависимый исполняемый файл

<h3 id="node">Нода</h3>

Список файлов в корне архива.

```bash
├── conf
├── contributors.txt
├── LICENSE
├── node-darwin-10.6-386
├── node-darwin-10.6-amd64
├── node-linux-386
├── node-linux-amd64
├── node-linux-arm-5
├── node-linux-arm-6
├── node-linux-arm64
├── node-linux-arm-7
├── node-linux-mips64
├── node-linux-mips64le
├── node-windows-4.0-386.exe
├── node-windows-4.0-amd64.exe
└── README.md
```

*   **conf** - директория конфигурационных файлов
*   **LICENSE** - текстовый файл лицензионного соглашения
*   **README.md** - текстовый файл краткой вводной информации
*   **contributors.txt** - текстовый файл, список разработчиков системы умный дом
*   **node-** - Платформо зависимый исполняемый файл
  

