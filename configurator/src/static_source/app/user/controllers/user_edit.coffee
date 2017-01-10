angular
.module('appControllers')
.controller 'userEditCtrl', ['$scope', 'Notify', 'User', '$stateParams', 'Message', '$state'
($scope, Notify, User, $stateParams, Message, $state) ->

  $scope.user = new User {id: $stateParams.id}

  show =->
    success =->
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
    return if confirm('точно удалить пользователя?')
    success =->
      $state.go 'dashboard.user.index'
    error =(result)->
      Message result.data.status, result.data.message
    $scope.user.$delete success, error

  show()

]