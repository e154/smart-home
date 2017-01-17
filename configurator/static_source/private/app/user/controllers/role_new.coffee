angular
.module('appControllers')
.controller 'roleNewCtrl', ['$scope', 'Notify', 'Role', '$stateParams', '$state', '$http', 'Message'
($scope, Notify, Role, $stateParams, $state, $http, Message) ->


  $scope.role = new Role {
    name: ''
    description: ''
  }

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
    success =(role)->
      Notify 'success', 'Роль успешно создана', 3
      $state.go 'dashboard.role.edit', {name: role.name}
    error =(result)->
      Message result.data.status, result.data.message
    $scope.role.$create success, error
]