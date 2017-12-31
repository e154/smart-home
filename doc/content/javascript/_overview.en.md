---
weight: 1
title: overview
groups:
    - javascript
---

## Обзор {#overview}

**API SmartHome**. Объекты и функции, методы, свойства объектов, используемых для управления 
устройствами "умного дома"

## Методы

Объект IC - глобальный объект, содержащий в себе все функции и константы API SmartHome.

<img src="/smart-home/img/schematic/workflow.svg" alt="smart-home workflow schematic map">

<img src="/smart-home/img/schematic/flow.svg" alt="smart-home flow schematic map">

**пример скрипта**

```coffeescript
# Контекст применения: 
# action (действие)
#
# Описание:
# Проверка состояния устройства. (частное)
# Не имеет зависимостей, и ни чего не передает наружу
# Должен вызываться в рамках воркера, или действия устройства,
# иначе выдаст ошибку, так как контекст выполнения накладывает 
# некоторые ограничения

fetchStatus =(node, dev)->

    # номер комманнды 
    # 3 - проверка состояния
    # 4 - выполнить комманду
    FUNCTION = 3

    # получим адрес устройства из контекста запуска
    DEVICE_ADDR = dev.getAddress()
    
    COMMAND = [DEVICE_ADDR, FUNCTION, 0, 0, 0, 5]
    
    # получить инстанс элемента на карте
    element = IC.Map.getElement dev
    
    # можно вывести произвольный текст под элементом
    # для отображения актуального состояния
    element.setOptions {text: 'some state'}
    
    # запрос состояния устройства
    from_node = node.send 'modbus', dev, true, COMMAND
    
    # запрос завершился c ошибкой    
    if from_node.error
        message.setError from_node.error
        
        # указать состояние, элемент автоматически изменит внешний вид
        # в зависимости от настроек состояний
        element.setState 'ERROR'
        
        # IC.Log.error "#{dev.name} - error: #{from_node.error}"
        return false
       
    # запрос отработал без ошибок, и что-то вернул
    if from_node.result != ""
    
        # так как тип запроса был modbus 
        # для работы с ответом нужно преобразовать его к массиву 
        result = IC.hex2arr(from_node.result)
        
        # в данном случае 1 означает что устройство включено
        # и функционирует, соответствунно сотоянию выставим 
        # состояние элемента карты
        if result[2] == 1
            element.setState 'ENABLED'
        else
            element.setState 'DISABLED'
    
    # print 'dev:', DEVICE_ADDR, 'state', element.getState().systemName
    
    from_node.result

main =->
    
    node = IC.CurrentNode()
    dev = IC.CurrentDevice()
    
    return if !node || !dev
    
    fetchStatus(node, dev)
    
```

<img src="/smart-home/img/schematic/screenshot.png" alt="smart-home scripts">

