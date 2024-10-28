---
title: "Execute"
linkTitle: "execute"
date: 2021-11-19
description: >

---

The **Smart Home** project's system provides the ability to execute arbitrary files and scripts synchronously and
asynchronously.

{{< alert color="success" >}}This function is available in any system script.{{< /alert >}}

To achieve this, the following methods are available:

1. `ExecuteSync(file, args)`: This method allows you to execute files and scripts synchronously. You pass the file name
   or path to the script as the `file` argument and any required arguments as an object `args`. Here's an example:

   ```javascript
   const file = 'script.js';
   const args = { param1: 'value1', param2: 'value2' };

   ExecuteSync(file, args);
   ```

2. `ExecuteAsync(file, args)`: This method allows you to execute files and scripts asynchronously. The `file` and `args`
   arguments have the same structure as in the `ExecuteSync` method. Here's an example:

   ```javascript
   const file = 'script.js';
   const args = { param1: 'value1', param2: 'value2' };

   ExecuteAsync(file, args);
   ```

Both the `ExecuteSync` and `ExecuteAsync` methods provide the ability to execute arbitrary files and scripts in your *
*Smart Home** project. The synchronous mode means that the code execution will block subsequent operations until the
script is completed, while the asynchronous mode allows you to continue executing other operations without waiting for
the script to finish. You can use these methods for integration with other systems, executing custom scripts, or
launching external applications in your **Smart Home** project.

----------------

### Code Example

```coffeescript
# ExecuteSync
# ##################################
"use strict";

r = ExecuteSync "data/scripts/ping.sh", "google.com"
if r.out == 'ok'
  print "site is available ^^"
```

ping.sh

```bash
#!/usr/bin/env bash

ping -c1 $1 > /dev/null 2> /dev/null; [[ $? -eq 0 ]] && echo ok || echo "err"
```
