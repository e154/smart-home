angular
.module('appServices')
.factory 'DeviceAction', ['$resource', ($resource) ->
  $resource window.server_url + '/api/v1/device_action/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.action || data

    create:
      method: 'POST'
      responseType: 'json'

    update:
        method: 'PUT'

    delete:
      method: 'DELETE'
]
