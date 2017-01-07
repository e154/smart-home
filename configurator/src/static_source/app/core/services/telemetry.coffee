angular
.module('appServices')
.factory 'Telemetry', ['$resource', ($resource) ->
  $resource window.server_url + '/api/v1/telemetry/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.telemetry || data

    create:
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        data?.telemetry || data

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
        items: data?.telemetrys || []
]
