#
# %div(map-viewer="options" map="map")
#

angular
.module('angular-map')
.directive 'mapViewer', ['$compile', ($compile) ->
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

    getOptions =->
      $scope.options = {} if !$scope.options
      options = $.extend true, defaultOptions, $scope.options

    $scope.$watch 'options', (val, oldVal)->
      return if !val || val == oldVal
      getOptions()


    #init
    getOptions()

    console.log 'options',$scope.options
    console.log 'map',$scope.map

    return
]