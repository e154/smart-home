angular
.module('angular-map')
.factory 'mapEditor', ['$rootScope', '$compile', 'mapFullscreen', 'mapPanning', '$templateCache', 'MapLayer', 'storage', '$timeout'
  ($rootScope, $compile, mapFullscreen, mapPanning, $templateCache, MapLayer, storage, $timeout) ->
    class mapEditor

      constructor: ()->
        @scope.addLayer = @addLayer
        @scope.removeLayer = @removeLayer
        @scope.updateLayer = @updateLayer
        @scope.selectLayer = @selectLayer
        @scope.selectElement = @selectElement
        @scope.addElement = @addElement
        @scope.current_layer = null
        @scope.sortLayers = @sortLayers
        @scope.sortElements = @sortElements
        @scope.removeElement = @removeElement
        @scope.updateElement = @updateElement
        @scope.addNewImage = @addNewImage
        @scope.addNewText = @addNewText
        @scope.addNewDevice = @addNewDevice
        @scope.addNewScript = @addNewScript
        @scope.preview = @preview
        @scope.current_element = {}

        @map_editor = new storage('map-editor')

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
        window_height = @map_editor.getInt('window_height')
        if window_height && window_height != 0
          @wrapper.height(window_height)

        if @wrapper.resizable('instance')
          @wrapper.resizable('destroy')
        @wrapper.resizable
          minHeight: @scope.settings.minHeight
          minWidth: @scope.settings.minWidth
          grid: @scope.settings.grid
          handles: 's'
          stop: (e)=>
            @map_editor.setItem('window_height', $(e.target).height())

        @panning = new mapPanning(@container, @scope, @wrapper)
        @wrapper.find(".page-loader").fadeOut("fast")

        @scope.$on 'select_element_on_map', (e, data)=>
          @selectElement(data)

        return

      addLayer: ()=>
        if !@model?.layers
          @model.layers = []

        layer = new MapLayer(@scope)
        layer.map_id = @id
        layer.create()
        @model.layers.unshift layer
        @selectLayer(@model.layers[@model.layers.length - 1])
        @sortLayers()

      removeLayer: (_layer)=>
        index = @model.layers.indexOf(_layer)
        success =()=>
          if index > -1
            @model.layers.splice(index, 1)
            @scope.current_layer = null

            if @model.layers.length > 0
              @selectLayer(@model.layers[@model.layers.length - 1])

        _layer.remove success
        return

      updateLayer: (_layer)=>
        return if !_layer
        _layer.update()

      selectLayer: (layer, $index)=>
        @scope.current_layer = layer
        angular.forEach @model.layers, (layer)=>
          angular.forEach layer.elements, (element)=>
            element.selected = false

      selectElement: (element)=>

        for layer in @model.layers
          if layer.id == element.layer_id
            @selectLayer(layer)
            break

        if @scope.current_element && @scope.current_element.id == element.id
          @scope.current_element = null
          element.selected = false
          $timeout ()=>
            @scope.$apply()
          return

        angular.forEach @model.layers, (layer)=>
          angular.forEach layer.elements, (_element)=>
            _element.selected = element.id == _element.id
            if _element.selected
              @scope.current_element = _element

        $timeout ()=>
          @scope.$apply()

      addElement: ()=>
        return if !@scope.current_layer
        element = @scope.current_layer.addElement()
        @selectElement(element)

      removeElement: (_element)=>
        return if !_element

        index = @scope.current_layer.elements.indexOf(_element)
        success =()=>
          if index > -1
            @scope.current_layer.elements.splice(index, 1)
            if @scope.current_layer.elements.length > 0
              @selectElement(@scope.current_layer.elements[@scope.current_layer.elements.length - 1])
            # for remove from scheme
            $timeout ()=>
              @scope.$apply()

        _element.remove success
        return

      updateElement: (_element)=>
        return if !_element
        success =()=>
          @scope.$broadcast 'entity_update', _element

        _element.update success

      sortLayers: ()=>
        weight = 0
        for layer in @model.layers
          layer.weight = weight
          weight++

      sortElements: ()=>
        return if !@scope.current_layer
        weight = 0
        for element in @scope.current_layer.elements
          element.weight = weight
          element.update_element_only()
          weight++

      addNewImage: ()=>
        return if !@scope.current_layer
        @scope.current_layer.addNewImage()

      addNewDevice: ()=>
        return if !@scope.current_layer
        @scope.current_layer.addNewDevice()

      addNewText: ()=>
        return if !@scope.current_layer
        @scope.current_layer.addNewText()

      addNewScript: ()=>
        return if !@scope.current_layer
        @scope.current_layer.addNewScript()

      preview: ()=>
        console.log 'preview'

    mapEditor
]