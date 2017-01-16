###* Copyright (C), DeltaSync Studios, 2010-2013. All rights reserved.
# ------------------------------------------------------------------
# File name:   editable.js
# Version:     v1.00
# Created:     20:07 / 19.06.14
# Author:      Delta54 <support@e154.ru>
#
# This file is part of the CMS engine (http://e154.ru/).
#
# Your use and or redistribution of this software in source and / or
# binary form, with or without modification, is subject to: (i) your
# ongoing acceptance of and compliance with the terms and conditions of
# the DeltaSync License Agreement; and (ii) your inclusion of this notice
# in any version of this software that you use or redistribute.
# A copy of the DeltaSync License Agreement is available by contacting
# DeltaSync Studios. at http://e154.ru/
#
# Description:
# ------------------------------------------------------------------
# History:
#
###

$.fn.editable = (options) ->
  addTa = (own) ->
    own.editing = true
    # отметка о начале редактирования
    own.html = $(own).html()
    # сохраним содержимое в буфере
    own.innerHTML = ''
    # сброс контейнер
    own.classList.add options.editClass
    # добавить класс означающий редактирование
    # замена всех <br> на \n
    content = own.html.replace(/(<br>)/g, '\n')
    # создадим форму ввода
    textarea = $('<textarea />')
    textarea.appendTo($(own)).attr('rows', getRows(content)).val(content).attr('cols', 20).attr('wrap',
      'soft').attr('name', $(own).attr('data-id')).attr('placeholder', $(own).attr('data-id')).prop('autofocus',
      true).prop('required', true).css('width': '100%').focus().focusout(->
      removeTa own, textarea
      return
    ).keyup(->
      @rows = getRows(@value)
      return
    ).keydown (e) ->
      if e.keyCode == 27
# ESC
        e.preventDefault()
        reset own, textarea
      return
    return

  removeTa = (own, t) ->
# Удаление формы
    own.editing = false
    # замена всех \r \n символов на <br>, и запись в буфер
    own.html = t.val().replace(/(\n|\r|\r\n)/g, '<br>')
    # удаление тега
    t.remove()
    # запись из буфера в контейнер
    # по дороге почистим от скриптов и прочей нечисти
    $(own).html own.html.replace(/<\s*script\s*.*>.*<\/\s*script\s*.*>/gi, '')
    # удалить класс options.editClass
    own.classList.remove options.editClass
    return

  arrayMerge = (a, b) ->
    if a
      if b
        for i of b
          a[i] = b[i]
      a
    else
      b

  getRows = (text) ->
# подсчёт строчек в форме
    getText = text
    getRegs = getText.match(/^.*(\r\n|\n|$)/gim)
    setText = false
    i = 0
    while i < getRegs.length
      getText = getRegs[i].replace(/\r|\n/g, '')
      setText += if getText.length then Math.ceil(getText.length / 50) else 1
      i++
    setText

  reset = (own, t) ->
# сброс и удаление формы
    t.remove()
    $(own).html own.html
    own.classList.remove options.editClass
    own.editing = false
    return

  options = arrayMerge({'editClass': 'editing'}, options)
  @click ->
    if @editing
      return
    addTa this
    return
  return
