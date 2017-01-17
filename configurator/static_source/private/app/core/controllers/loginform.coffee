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

  $scope.auth = ->
    success =(result)->
      authService.loginConfirmed()
      $scope.closeThisDialog()

    error = (result)->
      Message result.data.status, result.data.message
    $scope.user.$signin success, error

]
