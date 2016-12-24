angular
.module('appServices')
.factory 'MapLayer', ['$resource', ($resource) ->
  $resource window.server_url + '/api/v1/map_layer/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.map_layer || data

    create:
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        data?.map_layer || data

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
        items: data?.map_layers || []
]
