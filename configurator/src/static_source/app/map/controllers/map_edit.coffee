angular
.module('appControllers')
.controller 'mapEditCtrl', ['$scope', 'Map', '$state', 'Message', '$stateParams', 'mapConstructor', 'Notify'
($scope, Map, $state, Message, $stateParams, mapConstructor, Notify) ->

  $scope.map = new mapConstructor($scope, $stateParams.id)

  $scope.remove =->
    success =(data)->
      $state.go("dashboard.map.index")
    $scope.map.remove success

  $scope.update =->
    success =(data)->
      Notify 'success', 'Карта успешно обновлена', 3
    $scope.map.update success

  #------------------------------------------------------------------------------
  # init
  #------------------------------------------------------------------------------
  $scope.map.load()

  return
]