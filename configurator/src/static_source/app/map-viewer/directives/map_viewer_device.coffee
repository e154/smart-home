angular
.module('angular-map')
.directive 'mapViewerDevice', ['$compile', ($compile) ->
  restrict: 'A'
  replace: true
  templateUrl: '/map-viewer/templates/map_viewer_device.html'
  scope:
    element: '=mapViewerDevice'
  link: ($scope, $element, $attrs) ->
    $scope.element.click =()->
      console.log 'click'


    return
]