angular
.module('appServices')
.factory 'ImageResource', ['$resource', ($resource) ->
  $resource window.app_settings.server_url + '/api/v1/image/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.image || data

    create:
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        data?.image || data

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
        items: data?.images || []
]
