angular
.module('angular-map')
.directive 'mapEntity', ['$compile', '$templateCache', ($compile, $templateCache) ->
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

      # set params
      $element.css
        left: graph_settings.position.left || 0
        top: graph_settings.position.top || 0

      template = ''
      switch $scope.element.prototype_type
        when 'text'
          $element.css
            width: 'auto'
            height: 'auto'
          template = "<span class='draggable-entity' ng-style='{{element.prototype.style}}'>{{element.prototype.text}}</span>"
        when 'image'
          template = "<img class='draggable-entity' ng-src='{{element.prototype.image.url}}' ng-style='{{element.prototype.style}}'>"
          $element.css
            width: graph_settings.width || 'auto'
            height: graph_settings.height || 'auto'
        when 'device'
          template = $templateCache.get('/map-constructor/templates/_map_device_template.html')
          $element.css
            width: graph_settings.width || 'auto'
            height: graph_settings.height || 'auto'
        when 'script'
          break

      previousContent = $compile(template)($scope)
      $element.html(previousContent)

      if ['device', 'image'].indexOf($scope.element.prototype_type) > -1
        # set resizable
        if container.resizable('instance')
          container.resizable('destroy')
        container.resizable
          aspectRatio: true
          stop: (e)=>
            graph_settings.height = container.height()
            graph_settings.width = container.width()
            update()

    click =(e)->
      $scope.$emit 'select_element_on_map', $scope.element

    dragging =(e)->
      graph_settings.position.top = parseInt($(e.target).position().top, 10)
      graph_settings.position.left = parseInt($(e.target).position().left, 10)

    stop =(e)->
      update()

    update =->
      $scope.element.graph_settings = angular.copy(graph_settings)
      $scope.element.update_element_only()

    # --------------------
    # states:
    # enable/disable/frozen
    #

    isDraggable = false
    addDraggable =->
      if isDraggable
        container.draggable('enable')
        return
      isDraggable = true
      container.draggable(
        drag: (e)->
          dragging(e)
        stop: (e)->
          stop(e)
      ).on 'click', (e)->
        click()

    delDraggable =->
      container.draggable('disable')

    $scope.$watch 'element.status', (status)->
      return if !status || status == ''
      if status == 'frozen'
        delDraggable()
      else if status == 'enabled'
        addDraggable()

    # --------------------
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