
---
title: "Memory app"
linkTitle: "memory app"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/memory_app.png" >}}

&nbsp;

&nbsp;

In the **Smart Home** system, there is a "memory_app" plugin that provides the display of application memory parameters. This plugin allows you to obtain information about various characteristics of the application memory. Here are some of these parameters:

1. `alloc`: The `alloc` parameter displays the amount of memory allocated for the application in bytes. It indicates the current value of allocated memory.

2. `heap_alloc`: The `heap_alloc` parameter shows the amount of memory allocated in the heap for the application in bytes. This value indicates the current allocation of memory in the heap.

3. `total_alloc`: The `total_alloc` parameter indicates the total amount of memory allocated for the application since its launch in bytes. It reflects the overall memory volume allocated to the application.

4. `sys`: The `sys` parameter shows the amount of system memory used by the application. It indicates the current value of system memory usage.

5. `num_gc`: The `num_gc` parameter displays the number of garbage collections that have occurred since the application was launched.

6. `last_gc`: The `last_gc` parameter indicates the time of the last garbage collection. It provides information about the time of the last garbage collection performed in the application.

Here's an example of using the "memory_app" plugin to retrieve application memory parameters:

```javascript
const entity = entityManager.getEntity('memory.memory')
const memoryParams = entity.getAttributes()
console.log(memoryParams.alloc);
console.log(memoryParams.heap_alloc);
console.log(memoryParams.total_alloc);
console.log(memoryParams.sys);
console.log(memoryParams.num_gc);
console.log(memoryParams.last_gc);
```

You can use these parameters to display information about the application memory in your **Smart Home** project.
