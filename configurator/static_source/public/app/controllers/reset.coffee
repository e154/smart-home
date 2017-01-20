angular
.module('appControllers')
.controller 'resetCtrl', ['$scope', 'Auth'
($scope, Auth) ->

  $scope.email = ''
  $scope.error = null

  $scope.reset =->
    return if $scope.email == ''
    success =()->
      console.log 'ok'
    error =(error)->
      $scope.error = error
    Auth.reset {email: $scope.email}, success, error
]