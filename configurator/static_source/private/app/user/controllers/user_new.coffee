angular
.module('appControllers')
.controller 'userNewCtrl', ['$scope', 'Notify', 'User', 'Message', '$http', '$state'
($scope, Notify, User, Message, $http, $state) ->

  $scope.user = new User({
    status: 'active'
    role:
      name: 'user'
    meta: [
      {
        key: 'phone1'
        value: ''
      }
      {
        key: 'phone2'
        value: ''
      }
      {
        key: 'phone3'
        value: ''
      }
      {
        key: 'autograph'
        value: ''
      }
    ]

  })

  $scope.roles = []
  $scope.refreshRoles = (query)->
    $http(
      method: 'GET'
      url: window.app_settings.server_url + "/api/v1/role/search"
      params:
        query: query
        limit: 5
        offset: 0
    ).then (response)->
      $scope.roles = response.data.roles

  $scope.create =->
    success =(user)->
      Notify 'success', 'Пользователь успешно создан', 3
      $state.go 'dashboard.user.edit', {id: user.id}
    error =(result)->
      Message result.data.status, result.data.message
    $scope.user.$create success, error
]