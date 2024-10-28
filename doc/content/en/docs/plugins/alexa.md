---
title: "Alexa"
linkTitle: "alexa"
date: 2021-10-20
description: >
  
---

{{< alert color="warning" >}}Incomplete article{{< /alert >}}

In the **Smart Home** system, there is an Amazon Alexa plugin that provides integration with the Amazon Alexa voice
assistant. This plugin allows you to control devices and perform various operations in your smart home using voice
commands.

The Amazon Alexa plugin offers the following functionality:

1. Voice command recognition and processing: The plugin can receive voice commands from the user sent through the Amazon
   Alexa voice assistant. It recognizes the commands, analyzes them, and performs the corresponding actions in your
   smart home.

2. Device control: The plugin enables you to control devices in your smart home using voice commands. For example, you
   can say "Alexa, turn on the lights in the living room" or "Alexa, set the temperature to 25 degrees."

3. Integration with other system features: The Amazon Alexa plugin integrates with other features of the **Smart Home**
   system. For instance, you can use voice commands to trigger scenarios, control automated tasks, or obtain information
   about device status.

To use the Amazon Alexa plugin in your **Smart Home** project, you need to set it up and connect it to your Amazon Alexa
account. After that, you can configure voice commands and actions to control your smart home through the Amazon Alexa
voice assistant.

The Amazon Alexa plugin for the **Smart Home** system includes the following JavaScript handlers:

1. `skillOnLaunch`: This handler is called when the skill is launched on an Amazon Alexa device. It is intended to
   handle the initial skill launch event. You can define logic or perform necessary actions when the skill is launched.
   Example usage:

```javascript
skillOnLaunch = () => {
  // Logic for handling the skill launch event
};
```

2. `skillOnSessionEnd`: This handler is called when the session with the skill ends. It allows you to perform specific
   actions when the session with the user is terminated. For example, you can save state or perform final operations.
   Example usage:

```javascript
skillOnSessionEnd = () => {
  // Logic for session end
};
```

3. `skillOnIntent`: This handler is called when the skill receives intents from the user. It is intended to handle
   different intents and perform corresponding actions. You can define logic for intent handling and interaction with
   your smart home. Example usage:

```javascript
skillOnIntent = () => {
  // Logic for handling intents
};
```

These handlers provide the ability to handle skill launch events, session termination, and user intents within the
context of the Amazon Alexa plugin. You can define your own logic and perform necessary actions in each of these
handlers according to your project requirements.

### Code Example

```coffeescript
skillOnLaunch = ->
#print '---action onLaunch---'
  Done('skillOnLaunch')
skillOnSessionEnd = ->
#print '---action onSessionEnd---'
  Done('skillOnSessionEnd')
skillOnIntent = ->
#print '---action onIntent---'
  state = 'on'
  if Alexa.slots['state'] == 'off'
    state = 'off'

  place = Alexa.slots['place']

  Done("#{place}_#{state}")

  Alexa.sendMessage("#{place}_#{state}")
```
