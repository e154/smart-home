---
weight: 37
title: template
groups:
    - javascript
---

## IC.Template() {#ic_template}

Возвращает генератор шаблонов

### .render(name, options) {#ic_template_render}

Получить наименование устройства.

```coffeescript
tpl = IC.Template()
render = tpl.render(name, {'key':'val'})
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

