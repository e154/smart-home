angular
.module('angular-map')
.factory 'mapEditor', ['$rootScope', '$compile', 'mapFullscreen', 'mapPanning', '$templateCache', 'MapLayer', 'storage', '$timeout', 'MapLayerResource'
  ($rootScope, $compile, mapFullscreen, mapPanning, $templateCache, MapLayer, storage, $timeout, MapLayerResource) ->
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
        preventSelection(document.querySelector('.map-wrapper'))

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
          return if @scope.current_element && data.id == @scope.current_element.id
          @selectElement(data)

        return

      addLayer: ()=>
        if !@layers
          @layers = []

        layer = new MapLayer(@scope)
        layer.map_id = @id
        layer.create()
        @layers.unshift layer
        @selectLayer(layer)
        @sortLayers()

      removeLayer: (_layer)=>
        index = @layers.indexOf(_layer)
        success =()=>
          if index > -1
            @layers.splice(index, 1)
            @scope.current_layer = null

            if @layers.length > 0
              @selectLayer(@layers[0])

        _layer.remove success
        return

      updateLayer: (_layer)=>
        return if !_layer
        _layer.update()

      selectLayer: (layer, $index)=>
        @scope.current_layer = layer
#        @scope.current_element = null
        angular.forEach @layers, (layer)=>
          angular.forEach layer.elements, (element)=>
            element.selected = false

      selectElement: (element)=>

        for layer in @layers
          if layer.id == element.layer_id
            @selectLayer(layer)
            break

        if @scope.current_element && @scope.current_element.id == element.id
          @scope.current_element = null
          element.selected = false
          $timeout ()=>
            @scope.$apply()
          return

        angular.forEach @layers, (layer)=>
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

      removeSelected: ()=>
        return if !@scope.current_element
        @removeElement(@scope.current_element)

      removeElement: (_element)=>
        return if !@scope.current_layer
        cb =()=>
          index = @scope.current_layer.elements.indexOf(_element)
          if index > -1
            @scope.current_layer.elements.splice(index, 1)
            if @scope.current_layer.elements.length > 0
              @selectElement(@scope.current_layer.elements[0])
            else
              @scope.current_element = null
            # for remove from scheme
            $timeout ()=>
              @scope.$apply()
        @scope.current_layer.removeElement(_element, cb)

      updateElement: (_element)=>
        return if !_element
        success =()=>
          @scope.$broadcast 'entity_update', _element
        _element.update success

      sortLayers: ()=>
        weight = 0
        for layer in @layers
          layer.weight = weight
          weight++
        success =(data)->
        error =(result)->
          Message result.data.status, result.data.message
        MapLayerResource.sort @layers, success, error

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

      keyboard: ()=>
        angular.forEach @scope.settings.keyboard, (button, key_id)=>
          fn = this[button.callback] || window[button.callback]
          if typeof fn != 'function'
            return
          key key_id, (event, handler)=>
            fn.apply(null, [@scope])

    mapEditor
]