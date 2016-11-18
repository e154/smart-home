angular
.module('angular-bpmn')
.directive 'bpmnObject', ['log', '$timeout', '$templateCache', '$compile', '$templateRequest', '$q'
  (log, $timeout, $templateCache, $compile, $templateRequest, $q) ->
    restrict: 'A'
    controller: ["$scope", "$element", ($scope, $element)->

      container = $($element)
      childs = []

      switch $scope.data.type.name
        when 'poster'
          container.find('img').on 'dragstart', (e)->
            e.preventDefault()
        when 'group'
          if !($scope.data.resizable? && $scope.data.resizable) && $scope.object.settings.engine.status != 'editor'
            break

          container.resizable
            minHeight: 100
            minWidth: 100
            grid: 10
            handles: "all"
            start: (event, ui)->
              childs = $scope.object.getAllChilds()
            stop: (event, ui)->
              $scope.instance.repaintEverything()
            resize: (event, ui)->
              # во время изменения размера контейнера
              # контролирует нахлёст родительского блока с дочерними
              angular.forEach childs, (child, ch_key)->
                h = child.position.left + child.size.width
                v = child.position.top + child.size.height

                if container.width() <= h + 20
                  container.css('width', h + 20)

                if container.height() <= v + 20
                  container.css('height', v + 20)

              $scope.instance.repaintEverything()
        when 'swimlane'
          if !($scope.data.resizable? && $scope.data.resizable) && $scope.object.settings.engine.status != 'editor'
            break

          container.resizable
            minHeight: 200
            minWidth: 400
            grid: 10
            handles: 'e'
            start: ()->
              childs = $scope.object.getAllChilds()
              $scope.instance.repaintEverything()
            resize: ()->
              # во время изменения размера контейнера
              # контролирует нахлёст родительского блока с дочерними по левой стороне
              angular.forEach childs, (child)->
                if child.data.type == 'swimlane-row'
                  return
                h = child.position.left + child.size.width
                if container.width() <= h + 20
                  container.css('width', h + 20)
              $scope.instance.repaintEverything()
        when 'swimlane-row'
          if !($scope.data.resizable? && $scope.data.resizable) && $scope.object.settings.engine.status != 'editor'
            break

          container.resizable
            minHeight: 200
            minWidth: 400
            grid: 10
            handles: 's'
            start: ()->
              childs = $scope.object.getAllChilds()
              $scope.instance.repaintEverything()

            resize: ()->
              angular.forEach childs, (child)->
                v = child.position.top + child.size.height
                if container.height() <= v + 20
                  container.css('height', v + 20)

              $scope.instance.repaintEverything()

      updateStyle = ()->
        style =
          top: $scope.data.position.top
          left: $scope.data.position.left

        container.css(style)

        $scope.instance.repaintEverything()

      $scope.$watch 'data.position', (val, old_val) ->
        if val == old_val
          return
        updateStyle()

      updateStyle()
    ]
    link: ($scope, element, attrs) ->

]