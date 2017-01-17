angular
.module('appServices')
.factory 'Log', ['$resource', ($resource) ->
  $resource window.app_settings.server_url + '/api/v1/log/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.log || data

    create:
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        data?.log || data

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
        items: data?.logs || []
]
