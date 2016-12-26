angular
.module('appControllers')
.controller 'mapEditCtrl', ['$scope', '$state', 'Message', '$stateParams', 'mapConstructor', 'Notify'
($scope, $state, Message, $stateParams, mapConstructor, Notify) ->

  $scope.map = new mapConstructor($scope, parseInt($stateParams.id, 10))

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