
'use strict'

angular.module('templates', [])
angular.module('appDirectives', [])
angular.module('appFilters', [])
angular.module('appControllers', [])
angular.module('appServices', ['ngResource'])
angular.module('app', [
  'pascalprecht.translate'
  'templates'
  'appDirectives'
  'ngRoute'
  'appControllers'
  'appFilters'
  'appServices'
  'ui.router'
  'toaster'
])
