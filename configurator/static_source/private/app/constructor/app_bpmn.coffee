
'use strict'

angular.module('angular-bpmn', [])
angular.module('angular-bpmn')
.run ['$rootScope'
  ($rootScope) =>
    $rootScope.runMode = window.debug
]