---
title: "Javascript"
linkTitle: "Javascript"
date: 2021-11-19
description: >
  
---

The software logic and applications in the **Smart Home** project are developed using JavaScript, CoffeeScript, and
TypeScript. JavaScript is a widely-used programming language that provides flexibility and power in creating interactive
components.

The **Smart Home** project utilizes JavaScript to write all the core components that control the smart home. This
includes device management logic, data processing, interaction with external APIs, creating the user interface, and much
more.

In the **Smart Home** project, you can use JavaScript to create automation rules, schedule configurations, process data
from sensors, control lighting, heating, security systems, and many other aspects of your smart home. JavaScript enables
easy interaction with devices and systems and allows you to create user interfaces to control all the functions of your
home.

Thanks to the flexibility and popularity of JavaScript, you can easily customize and extend the **Smart Home** project
according to your needs and preferences. You have the opportunity to use a wide range of JavaScript tools and libraries
to adapt the project to your unique requirements and enhance its functionality.

Regardless of your programming skills, JavaScript provides a powerful toolkit for developing software logic and
applications in the **Smart Home** project.

### Scope and Visibility

{{< figure src="/smart-home/img/scrips-scope.svg" >}}

The script system for Smart Home provides a flexible and powerful environment for automation. However, it's essential to
consider different scopes when designing the software architecture. There are five scopes in which scripts can execute:

First, let's take a look at the big picture. We have five scopes where scripts can execute:

* **Entity**
    * **General Scope**
        * Executes once when the entity starts:
            * `function init(): void;`
        * Executes the `entityAction` method when an action is performed, if not specified in the **Action Scope**:
            * `function entityAction(entityId: string, actionName: string, params: { [key: string]: any }): void;`
    * **Action Scope**
        * Executes the `entityAction` method when an action is performed:
            * `function entityAction(entityId: string, actionName: string, params: { [key: string]: any }): void;`
* **Task**
    * **Trigger Scope**
        * Executes one of the following methods based on the trigger type and event:
            * `function automationTriggerStateChanged(msg: TriggerStateChangedMessage): boolean;`
            * `function automationTriggerAlexa(msg: TriggerAlexaTimeMessage): boolean;`
            * `function automationTriggerTime(msg: TriggerTimeMessage): boolean;`
            * `function automationTriggerSystem(msg: TriggerSystemMessage): boolean;`
    * **Condition Scope**
        * Executes the `automationCondition` method when the trigger is triggered:
            * `function automationCondition(entityId: string): boolean;`
    * **Action Scope**
        * Executes the `automationAction` method when `automationCondition` returns true:
            * `function automationAction(entityId: string): void;`

### Examples

We wrote a script describing logic, event handling, and automation.

```javascript
// Incorrect script implementation with three scopes

// General scope
// ##################################
var init, entityAction, automationTriggerStateChanged;

let entityState;

// Entity - General Scope
init = function () {
  // Entity initialization
  entityState = 'off';
}

// Action scope
// ##################################

entityAction = function (entityId, actionName, args) {
  console.log(entityState); // Incorrect
};

//  automation
//  ##################################
automationTriggerStateChanged = function (msg) {
  console.log(entityState); // Incorrect
};
```

```javascript
// Correct script implementation with three scopes

// General scope
// ##################################
var init, entityAction, automationTriggerStateChanged;

let entityState;
let storageKey = entityState + ENTITY_ID;

// Entity - General Scope
init = function () {
  // Entity initialization
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
