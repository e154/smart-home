angular
.module('appServices')
.factory 'Flow', ['$resource', ($resource) ->
  $resource window.server_url + '/api/v1/flow/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.flow || data

    create:
      method: 'POST'
      responseType: 'json'

    update:
        method: 'PUT'

    delete:
      method: 'DELETE'

    all:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        meta: data.meta
        items: data?.flows || []

]
