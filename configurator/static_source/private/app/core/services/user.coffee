angular
.module('appServices')
.factory 'User', ['$resource', ($resource) ->
  $resource window.app_settings.server_url + '/api/v1/user/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        user = data?.user || data
        user.full_name = "#{user.last_name} #{user.first_name}" if user.first_name && user.last_name || ""
        user

    create:
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        data?.user || data

    update_status:
      url: window.app_settings.server_url + '/api/v1/user/:id/update_status'
      method: 'PUT'
      responseType: 'json'

    update:
      method: 'PUT'
      responseType: 'json'

    delete:
      method: 'DELETE'

    all:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        meta: data?.meta || {}
        items: data?.users || []
]
