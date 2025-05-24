---
title: "Memory"
linkTitle: "memory"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/ram1.png" width="300" >}}

&nbsp;

&nbsp;

In the **Smart Home** system, there is a "memory" plugin that allows you to display the parameters of the RAM (random
access memory). This plugin provides access to the following parameters:

1. `total`: The `total` parameter displays the total amount of system memory in bytes.

2. `free`: The `free` parameter indicates the amount of free memory in bytes.

3. `used_percent`: The `used_percent` parameter shows the percentage of memory usage.

Here's an example of using the "memory" plugin to retrieve RAM parameters:

```javascript
const ramParams = EntityGetAttributes('memory.memory')
console.log(ramParams.total);
console.log(ramParams.free);
console.log(ramParams.used_percent);
```

You can use these parameters to display information about the RAM in your **Smart Home** project.
