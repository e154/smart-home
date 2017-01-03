angular
.module('angular-map')
.directive 'mapViewerElement', ['$compile', ($compile) ->
  restrict: 'A'
  templateUrl: '/map-viewer/templates/map_viewer_element.html'
  scope:
    element: '=mapViewerElement'
  link: ($scope, $element, $attrs) ->

    compile =->
      st = $scope.element.graph_settings

      $element.css
        left: st.position.left || 0
        top: st.position.top || 0

      switch $scope.element.prototype_type
        when 'text'
          $element.css
            width: 'auto'
            height: 'auto'
        else
          $element.css
            width: st.width || 'auto'
            height: st.height || 'auto'

    init =->
      $scope.$watch 'element.graph_settings', (val, oldVal)->
        return if !val || val == oldVal
        compile()

      compile()

    init()

    return
]