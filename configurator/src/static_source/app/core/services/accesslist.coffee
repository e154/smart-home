angular
.module('appServices')
.factory 'AccessList', ['$resource', ($resource) ->
  $resource window.server_url + '/api/v1/access_list/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.access_list || data

]
