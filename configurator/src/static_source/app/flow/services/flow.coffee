angular
.module('appServices')
.factory 'Flow', ['$resource', ($resource) ->
  $resource window.server_url + '/api/v1/flow/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.flow || data

    full:
      url: window.server_url + '/api/v1/flow/:id/full'
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.flow || data

    get_redactor:
      url: window.server_url + '/api/v1/flow/:id/redactor'
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.flow || data

    update_redactor:
      url: window.server_url + '/api/v1/flow/:id/redactor'
      method: 'PUT'

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
