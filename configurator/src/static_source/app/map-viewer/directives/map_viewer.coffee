#
# %div(map-viewer="options" map="map")
#

angular
.module('angular-map')
.directive 'mapViewer', ['$compile', 'mapPanning', 'mapFullscreen'
($compile, mapPanning, mapFullscreen) ->
  restrict: 'A'
  templateUrl: '/map-viewer/templates/map_viewer.html'
  scope:
    options: '=mapViewer'
    map: '='
  link: ($scope, $element, $attrs) ->

    # set default options
    options = {}
    defaultOptions =
      zoom: true

    $scope.zoom = 0.4
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

    getOptions =->
      $scope.options = {} if !$scope.options
      options = $.extend true, defaultOptions, $scope.options

    $scope.$watch 'options', (val, oldVal)->
      return if !val || val == oldVal
      getOptions()

    $scope.$watch 'map', (val, oldVal)->
      return if !$scope.map.layers
      angular.forEach $scope.map.layers, (layer)->
        angular.forEach layer.elements, (element)->
          element.graph_settings = angular.fromJson(element.graph_settings)
    , true

    #init
    getOptions()

#    console.log 'options',$scope.options
#    console.log 'map',$scope.map

    return
]