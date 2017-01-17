angular
.module('appServices')
.factory 'Auth', ['$resource', ($resource) ->
  $resource '/signin', {},
    signin:
      url: '/signin'
      method: 'POST'
      responseType: 'json'
      timeout: 60000

    recovery:
      url: '/recovery'
      method: 'POST'
      responseType: 'json'

    reset:
      url: '/reset'
      method: 'POST'
      responseType: 'json'
]
