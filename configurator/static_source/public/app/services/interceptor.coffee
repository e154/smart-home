angular
.module('appServices')
.factory('myHttpInterceptor', ['$q','$injector', '$rootScope', ($q, $injector, $rootScope) ->
  'request': (config) ->
    config.headers['X-Requested-With'] = 'XMLHttpRequest'
#    config.headers['Authorization'] = $rootScope.token if $rootScope.token
    config

  'requestError': (rejection) ->
    $q.reject rejection

  'response': (response) ->
    response

  'responseError': (rejection) ->
#    console.log "responseError ", rejection
    switch rejection.status
      when 401
        break
      when 403
        break
      when 404
        break
      when 500
        break
      else
        console.error "ERR_CONNECTION_REFUSED"

    $q.reject rejection

])

.config ['$httpProvider', ($httpProvider) ->
  $httpProvider.interceptors.push 'myHttpInterceptor'
]
