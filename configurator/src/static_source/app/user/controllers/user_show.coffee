angular
.module('appControllers')
.controller 'userShowCtrl', ['$scope', 'Notify','User', '$stateParams', '$state'
($scope, Notify, User, $stateParams, $state) ->

  $scope.user = new User {id: $stateParams.id}

  show =->
    success =->
    error =->
      $state.go 'dashboard.user.index'
    $scope.user.$show success, error

  show()

]