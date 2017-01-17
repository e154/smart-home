angular
.module('appServices')
.factory 'Device', ['$resource', ($resource) ->
  $resource window.app_settings.server_url + '/api/v1/device/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.device || data

    group:
      url: window.app_settings.server_url + '/api/v1/device/group'
      method: 'GET'
      responseType: 'json'

    actions:
      url: window.app_settings.server_url + '/api/v1/device/:id/actions'
      method: 'GET'
      isArray: true
      responseType: 'json'
      transformResponse: (data) ->
        data?.actions || data

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
        meta: data?.meta || {}
        items: data?.devices || []

]
