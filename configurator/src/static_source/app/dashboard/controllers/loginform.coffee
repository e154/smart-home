###*
# Created by delta54 on 09.11.14.
###

angular
.module('appControllers')
.controller 'loginFormCtrl', ['$scope','Auth','authService','Message', '$rootScope', 'storage', '$http'
($scope, Auth, authService, Message, $rootScope, storage, $http) ->

  $scope.user = new Auth {
    email: '',
    password: ''
  }

  auth = new storage('')
  $scope.auth = ->
    success =(result)->
#      TODO remove
      $rootScope.token = result.token
      $rootScope.current_user = result.current_user
      auth.setItem('token', result.token)
      auth.setObject('current_user', result.current_user)
      $http.defaults.headers.common['Authorization'] = $rootScope.token
      authService.loginConfirmed()
      $scope.closeThisDialog()

    error = (result)->
      Message result.data.status, result.data.message
    $scope.user.$signin success, error

]
