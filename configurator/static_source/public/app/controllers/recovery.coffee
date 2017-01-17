angular
.module('appControllers')
.controller 'recoveryCtrl', ['$scope', 'Auth'
($scope, Auth) ->

  $scope.password = ''
  $scope.email = ''
  $scope.error = null

  $scope.recovery =->
    success =()->
      console.log 'ok'
    error =(error)->
      $scope.error = error
    Auth.recovery {email: $scope.email, password: $scope.password}, success, error
]