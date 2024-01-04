
---
title: "Uptime"
linkTitle: "uptime"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/uptime1.png" width="300" >}}

&nbsp;

&nbsp;

In the **Smart Home** system, there is a plugin called "uptime" that allows displaying the parameters related to the system and application uptime. This plugin provides access to the following parameters:

1. `total`: The `total` parameter displays the total uptime of the system or application. It indicates the amount of time that has passed since the start.

2. `app_started`: The `app_started` parameter indicates the application's start time. It shows the date and time when the application was launched or restarted.

Here's an example of using the "uptime" plugin to retrieve uptime parameters:

```javascript
const uptimeParams = EntityGetAttributes('uptime.uptime')
console.log(uptimeParams.total);
console.log(uptimeParams.app_started);
```

You can use these parameters to display information about the system or application uptime in your **Smart Home** project.
