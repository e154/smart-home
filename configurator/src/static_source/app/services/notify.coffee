angular
.module("appServices")
.service 'Notify', ['$log', 'toaster'
($log, toaster) ->
  (status, message, time)->
    if !message || typeof message != 'string' || message == ""
      return

    if !time
      time = 5

    time *= 1000

    toaster.pop(status, null, message, time, 'trustedHtml')
]