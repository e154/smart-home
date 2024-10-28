---
title: "Sun"
linkTitle: "sun"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/sun1.png" width="300" >}}

&nbsp;

&nbsp;

The "sun" plugin allows you to retrieve information about the current position and time of the sun based on the
specified coordinates (latitude and longitude). It provides various parameters related to the position of the sun
throughout the day, as well as sunrise, sunset, twilight, and other moments. The plugin's state can indicate whether the
sun is above or below the horizon.

Parameters:

1. `horizonState`: Horizon state
2. `phase`: Sun phase
3. `azimuth`: Sun azimuth
4. `elevation`: Sun elevation above the horizon
5. `sunrise`: Sunrise time
6. `sunset`: Sunset time
7. `sunriseEnd`: Sunrise end time
8. `sunsetStart`: Sunset start time
9. `dawn`: Dawn time
10. `dusk`: Dusk time
11. `nauticalDawn`: Nautical dawn time
12. `nauticalDusk`: Nautical dusk time
13. `nightEnd`: Night end time
14. `night`: Night time
15. `goldenHourEnd`: Golden hour end time
16. `goldenHour`: Golden hour time
17. `solarNoon`: Solar noon time
18. `nadir`: Nadir (minimum sun elevation above the horizon)

Settings:

1. `lat`: Latitude (geographic latitude)
2. `lon`: Longitude (geographic longitude)

States:

1. `aboveHorizon`: The sun is above the horizon
2. `belowHorizon`: The sun is below the horizon
