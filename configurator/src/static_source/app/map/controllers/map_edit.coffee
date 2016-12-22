angular
.module('appControllers')
.controller 'mapEditCtrl', ['$scope', 'Map', '$state', 'Message', '$stateParams'
($scope, Map, $state, Message, $stateParams) ->

  $scope.map = {}

  #------------------------------------------------------------------------------
  # crud
  #------------------------------------------------------------------------------
  $scope.getMap =->
    success =(data)->
      $scope.map = angular.copy(data)
    error =(result)->
      Message result.data.status, result.data.message
    Map.show {id: $stateParams.id}, success, error

  $scope.remove =->
    return if !confirm('Вы точно хотите удалить эту карту?')
    success =(data)->
      $state.go("dashboard.map.index", {id: data.id})
    error =(result)->
      Message result.data.status, result.data.message
    $scope.map.$delete success, error

  $scope.update =->
    success =->
      Notify 'success', 'Карта успешно обновлена', 3
    error =(result)->
      Message result.data.status, result.data.message
    $scope.map.$update success, error


  #------------------------------------------------------------------------------
  # init
  #------------------------------------------------------------------------------
  $scope.getMap()

]