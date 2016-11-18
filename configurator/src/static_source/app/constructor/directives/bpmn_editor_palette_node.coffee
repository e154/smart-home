angular
.module('angular-bpmn')
.directive 'bpmnEditorPaletteNode', ['log', '$timeout', '$templateCache', '$compile', '$templateRequest'
  (log, $timeout, $templateCache, $compile, $templateRequest) ->
    restrict: 'A'
    scope:
      bpmnEditorPaletteNode: '='
      settings: '=settings'
    link: ($scope, element, attrs) ->
      container = $(element)
      data = $scope.bpmnEditorPaletteNode
      template = {}
      zoom = 1

      if data.shape.helper
        template = $compile('<div class="helper">'+data.shape.helper+'</div>')($scope)
      else if data.shape.templateUrl
        elementPromise = $templateRequest($scope.settings.theme.root_path + '/' + $scope.settings.engine.theme + '/' + data.shape.templateUrl)
        elementPromise.then (result)->
          template = $compile('<div class="helper" ng-class="[bpmnEditorPaletteNode.type.name]">'+result+'</div>')($scope)

      #TODO update zoom
      container.draggable({
        grid: $scope.settings.draggable.grid
        helper: ()->
          template.css({
            '-webkit-transform':zoom
            '-moz-transform':'scale('+zoom+')'
            '-ms-transform':'scale('+zoom+')'
            '-o-transform':'scale('+zoom+')'
            'transform':'scale('+zoom+')'
          })
#        appendTo: "body"
      })

      $scope.$parent.$parent.$parent.$parent.$watch 'zoom', (val, old_val)->
        zoom = val
]
