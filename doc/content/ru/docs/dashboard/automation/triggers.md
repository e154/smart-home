---
title: "Triggers"
linkTitle: "Triggers"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/triggers_window1.png" >}}

Триггеры в системе Smart Home представляют собой ключевые элементы, которые запускают выполнение Задач в ответ на
определенные события или условия. В настоящее время система поддерживает три основных типа триггеров, каждый из которых
обладает своими преимуществами и обеспечивает гибкость в настройке автоматизации.

#### 1. **Триггер "State Change" (Изменение Состояния):**

- **Описание:**
    - Активируется при изменении состояния определенного устройства или группы устройств в системе Smart Home.

- **Преимущества:**
    - **Реакция на Изменения:** Идеально подходит для сценариев, где требуется реагировать на конкретные изменения,
      например, включение/выключение света, открытие двери и т.д.
    - **Гибкость Настройки:** Позволяет выбирать конкретные устройства и параметры состояний для триггерирования.

#### 2. **Триггер "System" (Системные События):**

- **Описание:**
    - Активируется при возникновении системных событий, таких как запуск/остановка системы, подключение/отключение
      устройств и другие события системной шины.

- **Преимущества:**
    - **Мониторинг Системы:** Используется для мониторинга общего состояния системы и детектирования событий, связанных
      с функционированием умного дома.
    - **Интеграция с Внешними Системами:** Позволяет взаимодействовать с системами, работающими в рамках системной шины.

#### 3. **Триггер "Time (Cron)" (Временной Триггер):**

- **Описание:**
    - Активируется с заданным временным интервалом в соответствии с расписанием, заданным в формате cron.

- **Преимущества:**
    - **Регулярное Выполнение:** Идеально подходит для задач, которые должны выполняться регулярно по расписанию.
    - **Энергосбережение:** Может использоваться для оптимизации энергопотребления, например, выключения света или
      отопления в ночное время.

#### Примеры Использования Триггеров

1. **Триггер "State Change":**
    - **Сценарий:**
        - Включить кондиционер при открытии окна в спальне.
    - **Настройка Триггера:**
        - Устройство: Окно в спальне.
        - Состояние: Открыто.

2. **Триггер "System":**
    - **Сценарий:**
        - Отправить уведомление при подключении нового устройства к системе.
    - **Настройка Триггера:**
        - Системное событие: Подключение нового устройства.

3. **Триггер "Time (Cron)":**
    - **Сценарий:**
        - Выключить свет во всех комнатах после полуночи.
    - **Настройка Триггера:**
        - Временной интервал: "0 0 * * *", что соответствует каждой полуночи.

#### Примеры на Coffescript

1. `TriggerAlexa`:

```coffeescript
automationTriggerAlexa = (msg) ->
  p = msg.payload
  Done p
  return false
```

Обработчик `automationTriggerAlexa` вызывается в ответ на триггер от Amazon Alexa. Он принимает объект сообщения `msg`
и может выполнять определенные действия, связанные с этим триггером.

2. `TriggerStateChanged`:

```coffeescript
automationTriggerStateChanged = (msg) ->
  p = msg.payload
  Done p.new_state.state.name
  return false
```

Обработчик `automationTriggerStateChanged` вызывается при изменении состояния устройства. Он принимает объект сообщения
`msg` и может выполнить определенные действия на основе нового состояния устройства.

3. `TriggerSystem`:

```coffeescript
automationTriggerSystem = (msg) ->
  p = msg.payload
  Done p.event
  return false
```

Обработчик `automationTriggerSystem` вызывается в ответ на события системы. Он принимает объект сообщения `msg` и может
выполнить определенные действия, связанные с этим событием.

4. `TriggerTime`:

```coffeescript
automationTriggerTime = (msg) ->
  p = msg.payload
  Done p
  return false
```

Обработчик `automationTriggerTime` вызывается по истечении определенного времени. Он принимает объект сообщения `msg` и
может выполнить определенные действия, связанные с этим временем.

Каждый обработчик триггера может выполнять необходимую логику в ответ на соответствующий триггер, а затем возвращать
значение `false`, чтобы указать, что дальнейшая обработка не требуется.

Пример реализации:

```coffeescript
automationTriggerStateChanged = (msg)->
  p = msg.payload
  if !p.new_state || !p.new_state.state
    return false
  return msg.new_state.state.name == 'DOUBLE_CLICK'
```

#### Заключение

Использование различных типов триггеров в системе Smart Home позволяет пользователям создавать сложные и
интеллектуальные сценарии автоматизации, гармонично реагирующие на изменения в системе и внешние события. Это
предоставляет пользоват
