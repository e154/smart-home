angular
.module('appControllers')
.controller 'mapNewCtrl', ['$scope', 'Map', 'Message', '$state'
($scope, Map, Message, $state) ->

  $scope.map = new Map

  $scope.create =->
    success =(data)->
      $state.go("dashboard.map.edit", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    $scope.map.$create success, error

]