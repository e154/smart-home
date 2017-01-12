angular
.module('appControllers')
.controller 'roleShowCtrl', ['$scope', 'Notify','Role', '$stateParams', '$state'
($scope, Notify, Role, $stateParams, $state) ->

  $scope.role = new Role {name: $stateParams.name}

  show =->
    success =()->
    error =()->
      $state.go 'dashboard.role.index'
    $scope.role.$show success, error

  show()

]