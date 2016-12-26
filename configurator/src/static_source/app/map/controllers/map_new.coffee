angular
.module('appControllers')
.controller 'mapNewCtrl', ['$scope', 'MapResource', 'Message', '$state'
($scope, MapResource, Message, $state) ->

  $scope.map = new MapResource

  $scope.create =->
    success =(data)->
      $state.go("dashboard.map.edit.main", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    $scope.map.$create success, error

]