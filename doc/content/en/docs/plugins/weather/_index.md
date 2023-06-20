
---
title: "Weather"
linkTitle: "weather"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/weather1.png" >}}

&nbsp;

&nbsp;

In the **Smart Home** system, there is a basic "weather" plugin that provides the foundation for implementing weather plugins. This plugin enables access to weather information and allows the creation and integration of specific weather plugins for different data sources.

The "weather" plugin includes a list of weather conditions that can be used to display current weather conditions. Here is the list of weather conditions supported by the "weather" plugin:

1. `clearSky`: Clear sky.
2. `fair`: Fair weather.
3. `partlyCloudy`: Partly cloudy.
4. `cloudy`: Cloudy.
5. `lightRainShowers`: Light rain showers.
6. `rainShowers`: Rain showers.
7. `heavyRainShowers`: Heavy rain showers.
8. `lightRainShowersAndThunder`: Light rain showers and thunder.
9. `rainShowersAndThunder`: Rain showers and thunder.
10. `heavyRainShowersAndThunder`: Heavy rain showers and thunder.
11. `lightSleetShowers`: Light sleet showers.
12. `sleetShowers`: Sleet showers.
13. `heavySleetShowers`: Heavy sleet showers.
14. `lightSleetShowersAndThunder`: Light sleet showers and thunder.
15. `sleetShowersAndThunder`: Sleet showers and thunder.
16. `heavySleetShowersAndThunder`: Heavy sleet showers and thunder.
17. `lightSnowShowers`: Light snow showers.
18. `snowShowers`: Snow showers.
19. `heavySnowShowers`: Heavy snow showers.
20. `lightSnowShowersAndThunder`: Light snow showers and thunder.
21. `snowShowersAndThunder`: Snow showers and thunder.
22. `heavySnowShowersAndThunder`: Heavy snow showers and thunder.
23. `lightRain`: Light rain.
24. `rain`: Rain.
25. `heavyRain`: Heavy rain.
26. `lightRainAndThunder`: Light rain and thunder.
27. `rainAndThunder`: Rain and thunder.
28. `heavyRainAndThunder`: Heavy rain and thunder.
29. `lightSleet`: Light sleet.
30. `sleet`: Sleet.
31. `heavySleet`: Heavy sleet.

The "weather" plugin has basic settings that allow you to configure its behavior and the display of weather information. Here is a description of the basic settings of the "weather" plugin:

1. `lat`: Latitude - a numerical value indicating the geographic latitude of the location for retrieving weather data. Providing the correct latitude allows obtaining up-to-date weather for the specified region.
2. `lon`: Longitude - a numerical value indicating the geographic longitude of the location for retrieving weather data. Providing the correct longitude helps obtain up-to-date weather for the specified region.
3. `theme`: Theme - a parameter that determines the appearance and style of displaying weather data. Possible values may vary depending on the implementation of the plugin and user settings.
4. `winter`: Winter - a boolean value indicating which type of weather to display for the winter season. This setting allows adapting the display of weather depending on the time of year.

Please note that specific settings and their values may vary depending on the implementation of the "weather" plugin.

The "weather" plugin implements the following basic properties for displaying weather:

1. `main`: Main weather description (e.g., "Clear", "Cloudy", "Rain", etc.).
2. `datetime`: Date and time of the weather data.
3. `humidity`: Air humidity.
4. `max_temperature`: Maximum temperature.
5. `min_temperature`: Minimum temperature.
6. `pressure`: Atmospheric pressure.
7. `wind_bearing`: Wind direction.
   8

. `wind_speed`: Wind speed.

For each of the five days of weather forecast (day1, day2, day3, day4, day5), similar properties are available, with the corresponding day prefix. For example, for the first day of the forecast, properties like `day1_main`, `day1_datetime`, `day1_humidity`, etc., are accessible.

Example usage of the "weather" plugin's basic properties:

```javascript
const weatherPlugin = new WeatherPlugin();

// Get current weather data
const entity = entityManager.getEntity('weather_owm.weather1')
const currentWeather = entity.getAttributes()
console.log(currentWeather.main);
console.log(currentWeather.datetime);
console.log(currentWeather.humidity);
console.log(currentWeather.max_temperature);
console.log(currentWeather.min_temperature);
console.log(currentWeather.pressure);
console.log(currentWeather.wind_bearing);
console.log(currentWeather.wind_speed);

// Get weather data for the first day of the forecast
const entity = entityManager.getEntity('weather_owm.weather1')
const day1Weather = entity.getAttributes()
console.log(day1Weather.day1_main);
console.log(day1Weather.day1_datetime);
console.log(day1Weather.day1_humidity);
console.log(day1Weather.day1_max_temperature);
console.log(day1Weather.day1_min_temperature);
console.log(day1Weather.day1_pressure);
console.log(day1Weather.day1_wind_bearing);
console.log(day1Weather.day1_wind_speed);
```

Please note that specific properties and their values may vary depending on the implementation of the "weather" plugin and the weather data source.
