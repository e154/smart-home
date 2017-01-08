#
# %div(map-viewer="options" map="map")
#

angular
.module('angular-map')
.directive 'mapViewer', ['$compile', 'mapPanning', 'mapFullscreen', 'Stream', '$timeout'
($compile, mapPanning, mapFullscreen, Stream, $timeout) ->
  restrict: 'A'
  templateUrl: '/map-viewer/templates/map_viewer.html'
  scope:
    options: '=mapViewer'
    map: '='
  link: ($scope, $element, $attrs) ->

    # set default options
    # --------------------
    options = {}

    $scope.zoom = 1.0
    $scope.settings =
      movable: true
      zoom: true
      minHeight: 100
      minWidth: 400
    container = $element.find('.map-viewer')
    wrapper = $element.find('.map-wrapper')
    preventSelection(document.querySelector('.map-wrapper'))
    panning = new mapPanning(container, $scope, wrapper)
    fullscreen = new mapFullscreen(wrapper, $scope)

    $scope.devices = {}

    # stream
    # --------------------
    $timeout ()->
      Stream.sendRequest("get.devices.states", {}).then (data)->
        return if !data.states
        broadcastDeviceState(data.states)
    , 1000

    Stream.subscribe 'telemetry', (data)->
      return if !data.device
      state = {}
      state[data.device.id] = data.device.state
      broadcastDeviceState(state)

    broadcastDeviceState =(states)->
      angular.forEach states, (state, id)->
        id = parseInt(id, 10) if typeof id == 'string'
        $scope.$broadcast 'broadcast_device_state', {id: id, state: state}

    # etc
    # --------------------
    getOptions =->
      $scope.options = {} if !$scope.options
      options = $.extend true, $scope.settings, $scope.options

    $scope.$watch 'options', (val, oldVal)->
      return if !val || val == oldVal
      getOptions()

    $scope.$watch 'map.layers', (val, oldVal)->
      return if !val
      $scope.zoom = $scope.map.options?.zoom || 1.2
      panning.setZoom($scope.zoom)
      return if !$scope.map.layers
      angular.forEach $scope.map.layers, (layer)->
        angular.forEach layer.elements, (element)->
          element.graph_settings = angular.fromJson(element.graph_settings)

    #init
    getOptions()

#    console.log 'options',$scope.options
#    console.log 'map',$scope.map

    return
]