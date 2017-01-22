
'use strict'

### App Module ###

angular.module('templates', [])
angular.module('appDirectives', [])
angular.module('appFilters', [])
angular.module('appControllers', [])
angular.module('appServices', ['ngResource'])
app = angular.module('app', [
  'pascalprecht.translate'
  'templates'
  'appDirectives'
  'ngRoute'
  'appControllers'
  'appFilters'
  'appServices'
  'ui.router'
  'toaster'
  'bd.sockjs'
  'angular-bpmn'
  'ui.select'
  'ngSanitize'
  'ui.ace'
  'ngDialog'
  'ui.bootstrap.pagination'
  'pikaday'
  'ui.tree'
  'angular-map'
  'ng-sortable'
  'ngFileUpload'
  'ngDragDrop'
  'gridster'
  'passwordCheck'
  'http-auth-interceptor'
])

app.version = window.app_settings.version
