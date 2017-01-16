
'use strict'

angular.module('angular-map', [])
angular.module('angular-map')
.run ['$rootScope'
  ($rootScope) =>
    $rootScope.runMode = window.debug
]