angular
.module('appServices')
.factory 'Dashboard', ['$resource', ($resource) ->
  $resource window.server_url + '/api/v1/dashboard/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.dashboard || data

    create:
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        data?.dashboard || data

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
        items: data?.dashboards || []
]
