angular
.module('appControllers')
.controller 'dashboardIndexCtrl', ['$scope', 'Notify', 'Stream', 'Dashboard', '$timeout', 'Message'
($scope, Notify, Stream, Dashboard, $timeout, Message) ->

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

  # dashboard
  # --------------------
  $scope.dashboard = {}
  $scope.current_widget = null
  $scope.dashboard.widgets = [
    {id: 1, col: 0, row: 0, sizeY: 1, sizeX: 1, name: 'swap', type: 'swap' }
    {id: 2, col: 1, row: 0, sizeY: 1, sizeX: 2, name: 'memory', type: 'memory' }
    {id: 3, col: 3, row: 0, sizeY: 1, sizeX: 1, name: 'cpu', type: 'cpu_dig' }
    {id: 4, col: 0, row: 1, sizeY: 1, sizeX: 1, name: 'nodes', type: 'nodes' }
    {id: 5, col: 1, row: 1, sizeY: 1, sizeX: 1, name: 'devices', type: 'devices' }
  ]

  #TODO remove hard code
#  $scope.dashboard = new Dashboard({id:1})
#  success =(dashboard)->
#    dashboard.widgets = $scope.widgets
#  error = (result)->
#    Message result.data.status, result.data.message
#  $scope.dashboard.$show success, error

  # stream
  # --------------------
  $timeout ()->
    Stream.sendRequest("get.telemetry", {}).then (data)->
      return if !data.telemetry
      broadcastDeviceState(data.telemetry)
  , 1000

  Stream.subscribe 'telemetry', (data)->
    $scope.total_uptime = data.uptime.total if data.uptime?.total
    broadcastDeviceState(data)

  broadcastDeviceState =(data)->
    $scope.$broadcast 'telemetry_update', data

  $scope.$on '$stateChangeSuccess', ()->
    Stream.unsubscribe 'telemetry'

  # crud
  # --------------------
  $scope.removeWidget =(widget)->
    console.log 'remove widget', widget

  $scope.addWidget =(widget)->
    console.log 'add widget', widget

]