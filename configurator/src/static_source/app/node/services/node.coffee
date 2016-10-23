angular
.module('appServices')
.factory 'Node', ['$resource'
($resource) ->
  $resource window.server_url + '/api/v1/node/:id', {id: '@id'},
    get:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data.nodes

    post:
      method: 'POST'
      responseType: 'json'

    put:
        method: 'PUT'

    delete:
      method: 'DELETE'

    all:
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        meta: data.meta
        items: data.nodes

]
