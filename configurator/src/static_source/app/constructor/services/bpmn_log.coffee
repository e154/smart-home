angular
.module('angular-bpmn')
.service 'log', ['$log', '$rootScope'
  ($log, $rootScope) ->
    {
      debug: ()->
        if $rootScope.runMode != 'debug'
          return

        $log.debug.apply(self, arguments)

      error: ()->
        $log.error.apply(self, arguments)
    }
  ]