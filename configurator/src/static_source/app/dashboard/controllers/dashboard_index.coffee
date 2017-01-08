angular
.module('appControllers')
.controller 'dashboardIndexCtrl', ['$scope', 'Notify', 'Stream', 'Dashboard', '$timeout'
($scope, Notify, Stream, Dashboard, $timeout) ->

  angular.element(document).find('body').addClass('dashboard')
  $scope.total_uptime = 0

  # gridster
  # --------------------
  $scope.gridsterOptions =
    columns: 6,
    resizable:
      enabled: true,
      handles: ['e', 's', 'se'],
      stop: (event, $element, widget)-> {}
    draggable:
      stop: (event, $element, widget)-> {}

  $scope.dashboard = {
    id: 1
    name: 'Home'
    widgets: [
      {id: 1, col: 0, row: 0, sizeY: 1, sizeX: 1, name: 'swap', type: 'swap' }
      {id: 2, col: 1, row: 0, sizeY: 1, sizeX: 2, name: 'memory', type: 'memory' }
      {id: 3, col: 3, row: 0, sizeY: 1, sizeX: 1, name: 'cpu', type: 'cpu_dig' }
      {id: 4, col: 0, row: 1, sizeY: 1, sizeX: 1, name: 'nodes', type: 'nodes' }
      {id: 5, col: 1, row: 1, sizeY: 1, sizeX: 1, name: 'devices', type: 'devices' }
    ]
  }

  $timeout ()->
    Stream.sendRequest("get.telemetry", {}).then (data)->
      return if !data.telemetry
      broadcastDeviceState(data.telemetry)
  , 1000

  # remove
  # --------------------
  $scope.removeWidget =(widget)->
    console.log 'remove widget', widget

  # stream
  # --------------------
  Stream.subscribe 'telemetry', (data)->
    $scope.total_uptime = data.uptime.total if data.uptime?.total
    broadcastDeviceState(data)

  broadcastDeviceState =(data)->
    $scope.$broadcast 'telemetry_update', data

  $scope.$on '$stateChangeSuccess', ()->
    Stream.unsubscribe 'telemetry'

]