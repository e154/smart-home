angular
.module('appServices')
.factory 'Script', ['$resource', ($resource) ->
  $resource window.app_settings.server_url + '/api/v1/script/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.script || data

    create:
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        data?.script || data

    update:
        method: 'PUT'
        responseType: 'json'
        transformResponse: (data) ->
          data?.script || data

    delete:
      method: 'DELETE'

    all:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        meta: data?.meta || {}
        items: data?.scripts || []

    exec:
      url: window.app_settings.server_url + '/api/v1/script/:id/exec'
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        script = data?.script || data
        script.result = data?.result || ""
        script
]
