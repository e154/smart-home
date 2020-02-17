---
weight: 37
title: template
groups:
    - javascript
---

## Template {#template}

Возвращает генератор шаблонов

### .Render(name, options) {#template_render}

Получить наименование устройства.

```coffeescript
render = Template.Render(name, {'key':'val'})
```

**На входе**

**Значение** | **Описание**
-------------|--------------
  `name`     | type: string, наименование шаблона
  `options`  | type: Object, лъект для заполнения шаблона


**На выходе**

**Значение** | **Описание**
-------------|--------------
  `render`     | type: Object, шаблон

