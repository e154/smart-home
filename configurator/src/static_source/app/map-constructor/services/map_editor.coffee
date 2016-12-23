angular
.module('angular-map')
.factory 'mapEditor', ['$rootScope', '$compile', 'mapFullscreen', 'mapPanning', '$templateCache', 'mapLayer'
  ($rootScope, $compile, mapFullscreen, mapPanning, $templateCache, mapLayer) ->
    class mapElement

      constructor: ()->
        @scope.addLayer = @addLayer
        @scope.removeLayer = @removeLayer
        @scope.selectLayer = @selectLayer
        @scope.current_layer = null

      loadEditor: (c)=>
        # container
        # --------------------
        _container = angular.element(c)
        template = $templateCache.get('/map-constructor/templates/map_editor.html')
        angular.element(_container).append($compile(template)(@scope))

        @container = _container.find(".map-editor")
        @wrapper = _container.find(".map-wrapper")

        # fullscreen
        # --------------------
        @fullscreen = new mapFullscreen(@wrapper, @scope)

        # resizable
        # --------------------
        if @wrapper.resizable('instance')
          @wrapper.resizable('destroy')
        @wrapper.resizable
          minHeight: @scope.settings.minHeight
          minWidth: @scope.settings.minWidth
          grid: @scope.settings.grid
          handles: 's'

        @panning = new mapPanning(@container, @scope, @wrapper)
        @wrapper.find(".page-loader").fadeOut("fast")

        return

      addLayer: ()=>
        if !@model?.layers
          @model.layers = []

        @model.layers.push new mapLayer(@scope)

      removeLayer: (_layer)=>
        index = @model.layers.indexOf(_layer)
        if index > -1
          @model.layers.splice(index, 1)
          @scope.current_layer = null
        return

      selectLayer: (layer, $index)=>
        @scope.current_layer = layer


    mapElement
]