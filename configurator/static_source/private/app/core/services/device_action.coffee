angular
.module('appServices')
.factory 'DeviceAction', ['$resource', ($resource) ->
  $resource window.app_settings.server_url + '/api/v1/device_action/:id', {id: '@id'},
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

    get_by_device:
      url: window.app_settings.server_url + '/api/v1/device_action/get_by_device/:id'
      method: 'GET'
      isArray: true
      responseType: 'json'
      transformResponse: (data) ->
        data?.device_actions || data
]
