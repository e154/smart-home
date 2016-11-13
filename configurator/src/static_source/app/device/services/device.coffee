angular
.module('appServices')
.factory 'Device', ['$resource', ($resource) ->
  $resource window.server_url + '/api/v1/device/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.device || data

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
        items: data.devices

]
