angular
.module('appServices')
.factory 'Auth', ['$resource', ($resource) ->
  $resource '/signin', {},
    signin:
      url: '/signin'
      method: 'POST'
      responseType: 'json'
      timeout: 60000

    signout:
      url: '/signout'
      method: 'POST'
      responseType: 'json'

    show:
      url: window.app_settings.server_url + '/api/v1/access_list'
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.access_list || data

]
