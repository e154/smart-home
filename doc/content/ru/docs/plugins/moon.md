
---
title: "Moon"
linkTitle: "moon"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/moon1.png" width="300" >}}

&nbsp;

&nbsp;

В системе **Smart Home** присутствует плагин "moon", который позволяет отображать различные параметры, связанные с фазами луны и ее положением на небе. Этот плагин предоставляет доступ к следующим параметрам:

1. `horizonState`: Параметр `horizonState` указывает состояние луны относительно горизонта. Это может быть "aboveHorizon" (над горизонтом) или "belowHorizon" (под горизонтом).

2. `phase`: Параметр `phase` отображает текущую фазу луны. Например, "new_moon" (новолуние), "waxing_crescent" (растущая полумесяц), "first_quarter" (первая четверть), "waxing_gibbous" (растущая луна), "full_moon" (полнолуние), "waning_gibbous" (убывающая луна), "third_quarter" (третья четверть), "waning_crescent" (убывающая полумесяц).

3. `azimuth`: Параметр `azimuth` указывает азимутное положение луны в градусах.

4. `elevation`: Параметр `elevation` отображает угол возвышения луны над горизонтом.

5. `aboveHorizon`: Параметр `aboveHorizon` указывает, находится ли луна над горизонтом (true/false).

6. `belowHorizon`: Параметр `belowHorizon` указывает, находится ли луна под горизонтом (true/false).

Кроме того, плагин "moon" имеет настройки `lat` и `lon`, которые позволяют указать широту и долготу для определения положения луны на небе.

Пример использования плагина "moon" для получения параметров луны:

```javascript
const moonParams = EntityGetAttributes('moon.moon1')
console.log(moonParams.horizonState);
console.log(moonParams.phase);
console.log(moonParams.azimuth);
console.log(moonParams.elevation);
console.log(moonParams.aboveHorizon);
console.log(moonParams.belowHorizon);
```

Вы можете использовать эти параметры для отображения информации о луне в вашем проекте **Smart Home**.
