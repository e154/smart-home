angular
.module('appControllers')
.controller 'signinCtrl', ['$scope', 'Auth'
($scope, Auth) ->

  $scope.email = ''
  $scope.password = ''
  $scope.error = null

  $scope.signin =->
    success =()->
      window.location.reload()
    error =(error)->
      $scope.error = error
    Auth.signin {email: $scope.email, password: $scope.password}, success, error

]