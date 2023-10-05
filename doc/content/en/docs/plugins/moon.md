
---
title: "Moon"
linkTitle: "moon"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/moon1.png" >}}

&nbsp;

&nbsp;

In the **Smart Home** system, there is a "moon" plugin that allows displaying various parameters related to the phases of the moon and its position in the sky. This plugin provides access to the following parameters:

1. `horizonState`: The `horizonState` parameter indicates the moon's state relative to the horizon. It can be either "aboveHorizon" or "belowHorizon".

2. `phase`: The `phase` parameter displays the current phase of the moon. For example, "new_moon", "waxing_crescent", "first_quarter", "waxing_gibbous", "full_moon", "waning_gibbous", "third_quarter", "waning_crescent".

3. `azimuth`: The `azimuth` parameter indicates the azimuthal position of the moon in degrees.

4. `elevation`: The `elevation` parameter displays the angle of elevation of the moon above the horizon.

5. `aboveHorizon`: The `aboveHorizon` parameter indicates whether the moon is above the horizon (true/false).

6. `belowHorizon`: The `belowHorizon` parameter indicates whether the moon is below the horizon (true/false).

Additionally, the "moon" plugin has settings `lat` and `lon`, which allow specifying the latitude and longitude for determining the moon's position in the sky.

Here's an example of using the "moon" plugin to retrieve moon parameters:

```javascript
const moonParams = EntityGetAttributes('moon.moon1')
console.log(moonParams.horizonState);
console.log(moonParams.phase);
console.log(moonParams.azimuth);
console.log(moonParams.elevation);
console.log(moonParams.aboveHorizon);
console.log(moonParams.belowHorizon);
```

You can use these parameters to display information about the moon in your **Smart Home** project.
