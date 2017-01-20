angular
.module('appControllers')
.controller 'signinCtrl', ['$scope', 'Auth', '$timeout'
($scope, Auth, $timeout) ->

  $scope.email = ''
  $scope.password = ''
  $scope.error = null

  $scope.signin =->
    return if $scope.email == '' || $scope.password == ''
    success =()->
      window.location.reload()
    error =(error)->
      $scope.error = error.data.message
      $timeout ()->
        $scope.error = ''
      , 2000
    Auth.signin {email: $scope.email, password: $scope.password}, success, error

]