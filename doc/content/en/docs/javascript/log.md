---
title: "Log"
linkTitle: "log"
date: 2021-11-19
description: >

---

In the **Smart Home** project, there is a capability to log events using the "Log" object.

{{< alert color="success" >}}This function is available in any system script.{{< /alert >}}

The "Log" object provides the following methods for different logging levels:

1. `info(txt)`: This method is used to log informational messages. You pass the text message as the `txt` argument.
   Example usage:

```javascript
Log.info('This is an informational message.');
```

2. `warn(txt)`: This method is used to log warning messages. You pass the text message as the `txt` argument. Example
   usage:

```javascript
Log.warn('This is a warning message.');
```

3. `error(txt)`: This method is used to log error messages. You pass the error text message as the `txt` argument.
   Example usage:

```javascript
Log.error('An error occurred.');
```

4. `debug(txt)`: This method is used to log debug messages. You pass the debug text message as the `txt` argument.
   Example usage:

```javascript
Log.debug('Debugging information.');
```

The methods of the "Log" object allow you to log various types of messages such as informational messages, warnings,
errors, and debug messages. Logging helps track and analyze events occurring in the **Smart Home** project, making it
easier for the development, testing, and debugging process of your application.

----------------

### Code examples:

```coffeescript
# Log
# ##################################

Log.info 'some text'
Log.warn 'some text'
Log.error 'some text'
Log.debug 'some text'
```
