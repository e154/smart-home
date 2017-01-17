
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
  'templates'
  'appDirectives'
  'appControllers'
  'appFilters'
  'appServices'
  'ngRoute'
  'ui.router'
])

app.version =
  full: "#{APP_MAJOR}.#{APP_MINOR}.#{APP_PATCH}"
  major: APP_MAJOR
  minor: APP_MINOR
  patch: APP_PATCH
  developer: APP_DEVELOPER
  time: APP_RELEASE_TIME
