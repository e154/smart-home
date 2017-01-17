angular
.module('appServices')
.factory 'WidgetResource', ['$resource', ($resource) ->
  $resource window.app_settings.server_url + '/api/v1/widget/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.widget || data

    create:
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        data?.widget || data

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
        items: data?.widgets || []
]
