angular
.module('appServices')
.factory 'Auth', ['$resource', ($resource) ->
  $resource window.app_settings.server_url + '/api/v1/signin', {},
    signin:
      url: window.app_settings.server_url + '/api/v1/signin'
      method: 'POST'
      responseType: 'json'
      timeout: 60000
      transformResponse: (data) ->
        token = data?.token || null
        current_user = data?.current_user || null
        data

    signout:
      url: window.app_settings.server_url + '/api/v1/signout'
      method: 'POST'
      responseType: 'json'

    recovery:
      url: window.app_settings.server_url + '/api/v1/recovery'
      method: 'POST'
      responseType: 'json'

    reset:
      url: window.app_settings.server_url + '/api/v1/reset'
      method: 'POST'
      responseType: 'json'

    show:
      url: window.app_settings.server_url + '/api/v1/access_list'
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.access_list || data

]
