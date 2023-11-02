
---
title: "CPU"
linkTitle: "cpu"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/cpuspeed1.png" >}}

&nbsp;

&nbsp;

In the **Smart Home** system, there is a "cpuspeed" plugin that provides CPU parameter display. The plugin grants access to the following parameters:

1. `cores`: This parameter allows you to retrieve information about the number of CPU cores. It returns a number indicating the total number of CPU cores.

2. `mhz`: This parameter provides information about the current CPU frequency. It returns a value in megahertz (MHz) indicating the current operating frequency of the CPU.

3. `all`: The `all` parameter provides information about the weighted average CPU load in percentage.

4. `load_min`: This parameter displays the minimum CPU load over a specified period of time. It allows you to obtain information about the minimum CPU load and can be useful for monitoring and analyzing CPU load.

5. `load_max`: The `load_max` parameter provides information about the maximum CPU load over a specified period of time. It allows you to obtain information about the maximum CPU load and can be useful for analyzing system performance.

Here's an example of using the "cpuspeed" plugin to retrieve CPU parameters:

```javascript
const cpuParams = EntityGetAttributes('cpuspeed.cpuspeed')
console.log(cpuParams.cores);
console.log(cpuParams.mhz);
console.log(cpuParams.all);
console.log(cpuParams.load_min);
console.log(cpuParams.load_max);
```

You can use these parameters to display and monitor CPU information in your **Smart Home** project.
