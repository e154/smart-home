angular
.module('appControllers')
.controller 'userEditCtrl', ['$scope', 'Notify', 'User', '$stateParams', 'Message', '$state', '$http'
($scope, Notify, User, $stateParams, Message, $state, $http) ->

  $scope.user = new User {id: $stateParams.id}
  meta = [
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

  show =->
    success =(user)->
      if user.meta.length == 0
        $scope.user.meta = meta
    error =->
      $state.go 'dashboard.user.index'
    $scope.user.$show success, error

  $scope.removeAvatar =->
    $scope.user.avatar = null

  $scope.update =->
    success =->
      Notify 'success', 'Пользователь успешно обновлён', 3
    error =(result)->
      Message result.data.status, result.data.message
    $scope.user.$update success, error

  $scope.remove =->
    return if !confirm('точно удалить пользователя?')
    success =->
      $state.go 'dashboard.user.index'
    error =(result)->
      Message result.data.status, result.data.message
    $scope.user.$delete success, error

  show()

]