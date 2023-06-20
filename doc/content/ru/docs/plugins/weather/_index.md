
---
title: "Weather"
linkTitle: "weather"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/weather1.png" >}}

&nbsp;

&nbsp;

В системе **Smart Home** присутствует базовый плагин "weather", который предоставляет основу для реализации погодных плагинов. 
Этот плагин обеспечивает доступ к информации о погоде, позволяя создавать и интегрировать специфичные погодные плагины для различных источников данных.

Плагин "weather" включает список состояний погоды, которые могут быть использованы для отображения текущих условий погоды. 
Вот список состояний погоды, которые поддерживаются плагином "weather":

1. `clearSky`: Чистое небо.
2. `fair`: Ясно.
3. `partlyCloudy`: Частичная облачность.
4. `cloudy`: Облачно.
5. `lightRainShowers`: Небольшие кратковременные дожди.
6. `rainShowers`: Кратковременные дожди.
7. `heavyRainShowers`: Сильные кратковременные дожди.
8. `lightRainShowersAndThunder`: Небольшие кратковременные дожди и гроза.
9. `rainShowersAndThunder`: Кратковременные дожди и гроза.
10. `heavyRainShowersAndThunder`: Сильные кратковременные дожди и гроза.
11. `lightSleetShowers`: Небольшие кратковременные метель и дождь.
12. `sleetShowers`: Кратковременные метель и дождь.
13. `heavySleetShowers`: Сильные кратковременные метель и дождь.
14. `lightSleetShowersAndThunder`: Небольшие кратковременные метель, дождь и гроза.
15. `sleetShowersAndThunder`: Кратковременные метель, дождь и гроза.
16. `heavySleetShowersAndThunder`: Сильные кратковременные метель, дождь и гроза.
17. `lightSnowShowers`: Небольшие кратковременные снегопады.
18. `snowShowers`: Кратковременные снегопады.
19. `heavySnowShowers`: Сильные кратковременные снегопады.
20. `lightSnowShowersAndThunder`: Небольшие кратковременные снегопады и гроза.
21. `snowShowersAndThunder`: Кратковременные снегопады и гроза.
22. `heavySnowShowersAndThunder`: Сильные кратковременные снегопады и гроза.
23. `lightRain`: Небольшой дождь.
24. `rain`: Дождь.
25. `heavyRain`: Сильный дождь.
26. `lightRainAndThunder`: Небольшой дождь и гроза.
27. `rainAndThunder`: Дождь и гроза.
28. `heavyRainAndThunder`: Сильный дождь и гроза.
29. `lightSleet`: Небольшая метель и дождь.
30. `sleet`: Метель и дождь.
31. `heavySleet`: Сильная метель и дож

Плагин "weather" имеет базовые настройки, которые позволяют настроить его поведение и отображение погодной информации. 
Вот описание базовых настроек плагина "weather":

1. `lat`: Широта (latitude) - числовое значение, указывающее географическую широту местоположения для получения погодных 
данных. Задание правильной широты позволяет получить актуальную погоду для указанного региона.

2. `lon`: Долгота (longitude) - числовое значение, указывающее географическую долготу местоположения для получения погодных данных. 
Задание правильной долготы помогает получить актуальную погоду для указанного региона.

3. `theme`: Тема (theme) - параметр, определяющий внешний вид и стиль отображения погодных данных. Возможные значения 
могут варьироваться в зависимости от реализации плагина и пользовательских настроек.

4. `winter`: Зима (winter) - логическое значение, указывающее, какой вид погоды следует отображать для зимнего сезона. 
Эта настройка позволяет адаптировать отображение погоды в зависимости от времени года.

Пожалуйста, имейте в виду, что конкретные настройки и их значения могут отличаться в зависимости от реализации плагина "weather".

Плагин "weather" реализует следующие базовые свойства для отображения погоды:

1. `main`: Основное описание погоды (например, "Ясно", "Облачно", "Дождь" и т.д.).
2. `datetime`: Дата и время отображения погодных данных.
3. `humidity`: Влажность воздуха.
4. `max_temperature`: Максимальная температура.
5. `min_temperature`: Минимальная температура.
6. `pressure`: Атмосферное давление.
7. `wind_bearing`: Направление ветра.
8. `wind_speed`: Скорость ветра.

Для каждого из пяти дней прогноза погоды (day1, day2, day3, day4, day5) также доступны аналогичные свойства, с добавлением 
префикса соответствующего дня. Например, для первого дня прогноза доступны свойства `day1_main`, `day1_datetime`, `day1_humidity` и т.д.

Пример использования базовых свойств плагина "weather":

```javascript
const weatherPlugin = new WeatherPlugin();

// Получение данных о текущей погоде
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

// Получение данных о погоде на первый день прогноза
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

Обратите внимание, что конкретные свойства и их значения могут отличаться в зависимости от реализации плагина "weather" и источника погодных данных.
