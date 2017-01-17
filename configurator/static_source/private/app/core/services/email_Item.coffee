angular
.module('appServices')
.factory 'EmailItem', ['$resource', ($resource) ->
  $resource window.app_settings.server_url + '/api/v1/email/item/:name', {id: '@name'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.item || data  

    create:
      method: 'POST'
      responseType: 'json'

    update:
      method: 'PUT'

    delete:
      method: 'DELETE'

    items:
      url: window.app_settings.server_url + '/api/v1/email/items'
      method: 'GET'

    get_tree:
      url: window.app_settings.server_url + '/api/v1/email/items/tree'
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.tree || data

    update_tree:
      url: window.app_settings.server_url + '/api/v1/email/items/tree'
      method: 'POST'

    all:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        meta: data?.meta || {}
        items: data?.items || []

]
