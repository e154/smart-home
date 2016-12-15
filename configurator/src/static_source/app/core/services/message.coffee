angular
.module('appServices')
.factory 'Message', ['toaster','$filter','$translatePartialLoader','$translate'
(toaster, $filter, $translatePartialLoader, $translate) ->
  $translatePartialLoader.addPart 'messages'
  $translate.refresh()

  translateMessage = (msg, status) ->
    arr = msg.split(': ', 2)
    title = undefined
    body = undefined
    if arr.length > 1
      title = $filter('translate')(arr[0])
      body = $filter('translate')(arr[1])
    else
      title = $filter('translate')(status)
      body = $filter('translate')(arr[0])
    title + ': ' + body

  (status, message, time) ->
    if typeof message == 'undefined'
      return
    if typeof time == 'undefined'
      time = 5000
    arr = message.split('\r')
    if arr.length > 1
      i = 0
      while i < arr.length
        if arr[i] == ''
          i++
          continue
        toaster.pop status, null, translateMessage(arr[i], status), time, 'trustedHtml'
        i++
    else
      toaster.pop status, null, translateMessage(message, status), time, 'trustedHtml'

]
