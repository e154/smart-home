angular
.module('appControllers')
.controller 'dashboardIndexCtrl', ['$scope', 'Notify', 'Stream', 'Telemetry', 'Dashboard', 'Message'
($scope, Notify, Stream, Telemetry, Dashboard, Message) ->

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
#      {id: 4, col: 4, row: 0, sizeY: 1, sizeX: 1, name: 'uptime', type: 'uptime' }
    ]
  }

  success =(telemetry)->
    broadcastDeviceState(telemetry)
  error =(result)->
    Message result.data.status, result.data.message
  Telemetry.show {id: 1}, success, error

  # remove
  # --------------------
  $scope.removeWidget =(widget)->
    console.log 'remove widget', widget

  # stream
  # --------------------
  Stream.subscribe 'telemetry', (data)->
    $scope.total_uptime = data.uptime.total if data.uptime.total
    broadcastDeviceState(data)

  broadcastDeviceState =(data)->
    $scope.$broadcast 'telemetry_update', data

  $scope.$on '$stateChangeSuccess', ()->
    Stream.unsubscribe 'telemetry'

]