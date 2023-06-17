
---
title: "Triggers"
linkTitle: "Triggers"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/trigger_window1.png" >}}

&nbsp;

&nbsp;

Каждый тип триггера имеет свой собственный обработчик (handler) для выполнения соответствующих действий.

Вот примеры различных типов триггеров и соответствующих им обработчиков:

1. `TriggerAlexa`:
```coffeescript
automationTriggerAlexa = (msg) ->
  p = unmarshal msg.payload
  Done p
  return false
```
Обработчик `automationTriggerAlexa` вызывается в ответ на триггер от Amazon Alexa. Он принимает объект сообщения `msg`
и может выполнять определенные действия, связанные с этим триггером.

2. `TriggerStateChanged`:
```coffeescript
automationTriggerStateChanged = (msg) ->
  print '---trigger---'
  p = unmarshal msg.payload
  Done p.new_state.state.name
  return false
```
Обработчик `automationTriggerStateChanged` вызывается при изменении состояния устройства. Он принимает объект сообщения
`msg` и может выполнить определенные действия на основе нового состояния устройства.

3. `TriggerSystem`:
```coffeescript
automationTriggerSystem = (msg) ->
  p = unmarshal msg.payload
  Done p.event
  return false
```
Обработчик `automationTriggerSystem` вызывается в ответ на события системы. Он принимает объект сообщения `msg` и может
выполнить определенные действия, связанные с этим событием.

4. `TriggerTime`:
```coffeescript
automationTriggerTime = (msg) ->
  p = unmarshal msg.payload
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
    #print '---trigger---'
    p = unmarshal msg.payload
    if !p.new_state || !p.new_state.state
        return false
    return msg.new_state.state.name == 'DOUBLE_CLICK'
```
