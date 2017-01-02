angular
.module('appControllers')
.controller 'mapShowCtrl', ['$scope', '$state', 'Message', '$stateParams', 'MapResource', 'Notify'
($scope, $state, Message, $stateParams, MapResource, Notify) ->

  $scope.map = new MapResource {id: $stateParams.id}
  $scope.options = {}

  getMap =->
    success =()=>
    error =(result)->
      Message result.data.status, result.data.message
    $scope.map.$showFull success, error

  #------------------------------------------------------------------------------
  # init
  #------------------------------------------------------------------------------
  getMap()

  return
]