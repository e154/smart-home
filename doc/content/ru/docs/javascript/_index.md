
---
title: "Javascript"
linkTitle: "Javascript"
date: 2021-11-19
description: >
  
---

Программная логика и приложения в проекте **Smart Home** разрабатываются с использованием языка JavaScript, CoffeeScript и TypeScript. 
JavaScript - это широко используемый язык программирования, который обеспечивает гибкость и мощность в создании 
интерактивных компонентоа.

Проект **Smart Home** использует JavaScript для написания всех основных компонентов, которые управляют умным домом.
Это включает логику управления устройствами, обработку данных, взаимодействие с внешними API, создание пользовательского 
интерфейса и многое другое.

В проекте **Smart Home** вы можете использовать JavaScript для создания правил автоматизации, настройки расписаний, 
обработки данных с датчиков, управления освещением, отоплением, системой безопасности и многими другими аспектами вашего
умного дома. JavaScript обеспечивает простоту взаимодействия с устройствами и системами, а также позволяет создавать
пользовательские интерфейсы для управления всеми функциями вашего дома.

Благодаря гибкости и популярности JavaScript, вы можете легко настроить и расширить проект **Smart Home** в соответствии 
с вашими потребностями и предпочтениями. Вам предоставляется возможность использовать широкий спектр инструментов и 
библиотек JavaScript, чтобы адаптировать проект под ваши уникальные требования и расширить его функциональность.

Независимо от ваших навыков в программировании, JavaScript предоставляет мощный инструментарий для разработки программной
логики и приложений в проекте **Smart Home**.

### Область видимости и scope

{{< figure src="/smart-home/img/scrips-scope.svg" >}}

Система скриптов для Smart Home предоставляет гибкую и мощную среду для автоматизации, однако важно учесть различные
области видимости при проектировании программной архитектуры. Существует пять скоупов, в которых могут выполняться
скрипты:

Для начала давайте посмотрим картину в целом. У нас есть пять скоупов где исполняются скрипты:
* **Entity**
  * **General Scope**
    * при старте entity выполняется единожды:
      * `function init(): void;`
    * выполняется метод entityAction, когда, выполняется действие, если не указан в **Action Scope**:
      * `function entityAction(entityId: string, actionName: string, params: { [key: string]: any }): void;`
  * **Action Scope**
    * выполняется метод entityAction, когда, выполняется действие:
      * `function entityAction(entityId: string, actionName: string, params: { [key: string]: any }): void;`
* **Task**
  * **Trigger Scope**
    * выполняется один из следующих методов в зависимости от типа триггера и события:
      * `function automationTriggerStateChanged(msg: TriggerStateChangedMessage): boolean;`
      * `function automationTriggerAlexa(msg: TriggerAlexaTimeMessage): boolean;`
      * `function automationTriggerTime(msg: TriggerTimeMessage): boolean;`
      * `function automationTriggerSystem(msg: TriggerSystemMessage): boolean;`
  * **Condition Scope**
    * выполняется метод automationCondition, когда отработает триггер:
      * `function automationCondition(entityId: string): boolean;`
  * **Action Scope**
    * выполняется метод automationAction, когда automationCondition вернет true:
      * `function automationAction(entityId: string): void;`

### Примеры

Мы написали скрипт описывающий логику, обработчик событий и автоматизацию. 

```javascript
// Неправильная реализация скрипта с тремя скопами

// General scope
// ##################################
var init, entityAction, automationTriggerStateChanged;

let entityState;

// Entity - General Scope
init = function () {
  // Инициализация сущности
  entityState = 'off';
}

// Action scope
// ##################################

entityAction = function (entityId, actionName, args) {
  console.log(entityState); // Неправильно
};

//  automation
//  ##################################
automationTriggerStateChanged = function (msg) {
  console.log(entityState); // Неправильно
};
```

```javascript
// Правильная реализация скрипта с тремя скопами

// General scope
// ##################################
var init, entityAction, automationTriggerStateChanged;

let entityState;
let storageKey = entityState + ENTITY_ID;

// Entity - General Scope
init = function () {
  // Инициализация сущности
  entityState = 'off';

  Storage.push(storageKey, entityState);
}

// Action scope
// ##################################

entityAction = function (entityId, actionName, args) {
  const entityState = Storage.getByName(storageKey);
  console.log(entityState);
};

//  automation
//  ##################################
automationTriggerStateChanged = function (msg) {
  const entityState = Storage.getByName(storageKey);
  console.log(entityState);
};
```
