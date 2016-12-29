angular
.module('angular-map')
.directive 'mapEntity', ['$compile', ($compile) ->
  restrict: 'EA'
  replace: true
  scope:
    element: '=mapEntity'
  link: ($scope, $element, $attrs) ->

    previousContent = null
    container = $($element)
    graph_settings =
      width: null
      height: null
      position:
        top: 0
        left: 0

    compile =->
      if previousContent
        previousContent.remove()
        previousContent = null

      template = ''
      switch $scope.element.prototype_type
        when 'text'
        when 'image'
          template = "<img class='draggable-entity' ng-src=\"{{element.prototype.image.url}}\">"
        when 'device'
          break
        when 'script'
          break

      previousContent = $compile(template)($scope)
      $element.append(previousContent)

      # set params
      $element.css
        left: graph_settings.position.left || 0
        top: graph_settings.position.top || 0
        width: graph_settings.width || 'auto'
        height: graph_settings.height || '100px'

      # set resizable
      if container.resizable('instance')
        container.resizable('destroy')
      container.resizable
        aspectRatio: true
        stop: (e)=>
          graph_settings.height = container.height()
          graph_settings.width = container.width()
          update()

    addDraggable =->
      container.draggable(
        drag: (e)->
          dragging(e)
        stop: (e)->
          stop(e)
      ).on 'click', (e)->
        click()

    click =(e)->
      $scope.$emit 'select_element', $scope.element

    dragging =(e)->
      graph_settings.position.top = parseInt($(e.target).position().top, 10)
      graph_settings.position.left = parseInt($(e.target).position().left, 10)

    stop =(e)->
      update()

    update =->
      $scope.element.graph_settings = angular.copy(graph_settings)
      $scope.element.update_element_only()

    # init
    #
    graph_settings = angular.copy($scope.element.graph_settings)
    compile()
    addDraggable()

    $scope.$on 'entity_update', (e, data)->
      if data.id == $scope.element.id
        compile()

    return
]