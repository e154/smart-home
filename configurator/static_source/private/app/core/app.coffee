
'use strict'

APP_MAJOR = 1
APP_MINOR = 0
APP_PATCH = 0
APP_RELEASE_TIME = '__CURRENT_TIME__'
APP_DEVELOPER = 'delta54<support@e154.ru>'

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

app.version =
  full: "#{APP_MAJOR}.#{APP_MINOR}.#{APP_PATCH}"
  major: APP_MAJOR
  minor: APP_MINOR
  patch: APP_PATCH
  developer: APP_DEVELOPER
  time: APP_RELEASE_TIME
