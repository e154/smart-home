angular
.module('angular-map')
.factory 'MapLayer', ['$rootScope', '$compile', 'MapElement', 'MapLayerResource', 'Notify', 'Message', 'FileManager', 'MapElementResource'
  ($rootScope, $compile, MapElement, MapLayerResource, Notify, Message, FileManager, MapElementResource) ->
    class MapLayer

      scope: null

      id: null
      map_id: null
      name: 'Новый слой'
      description: ''
      status: 'enabled'
      elements: null
      created_at: null
      update_at: null
      weight: 0

      constructor: (@scope)->
        @elements = []

      serialize: ()->
        elements = []
        angular.forEach @elements, (element)->
          elements.push element.serialize()

        name: @name
        id: @id if @id
        map: {id: @map_id} if @map_id
        status: @status
        description: @description
        created_at: @created_at if @created_at
        update_at: @update_at if @update_at
        elements: elements if elements.length
        weight: @weight

      deserialize: (layer)->
        @id = layer.id || null
        @map_id = layer.map.id || null
        @name = layer.name || ''
        @description = layer.description || ''
        @status = layer.status || ''
        @created_at = layer.created_at || ''
        @update_at = layer.update_at || ''
        @weight = layer.weight || 0

        angular.forEach layer.elements, (element)=>
          @elements.push new MapElement(@scope).deserialize(element)

        return @

      addElement: ()=>
        element = new MapElement(@scope)
        element.layer_id = @id
        element.map_id = @map_id
        element.name += " №#{@elements.length + 1}"
        element.inheritPosition()
        element.create()
        @addToLayer(element)

        element

      copyElement: (_element, $event)=>
        $event.stopPropagation()
        $event.preventDefault()
        element = new MapElement(@scope)
        element.copy(_element)
        element.create()
        @addToLayer(element)

      removeElement: (_element, cb)->
        return if !_element || $.isEmptyObject(_element)
        success =()=>
          cb() if cb
        _element.remove success
        return

      create: ()->
        success =(data)=>
          @id = data.id
          Notify 'success', 'Слой успешно создан', 3
        error =(result)->
          Message result.data.status, result.data.message

        model = new MapLayerResource(@serialize())
        model.$create success, error

      update: (cb)->
        success =(data)=>
          Notify 'success', 'Слой успешно обновлён', 3
        error =(result)->
          Message result.data.status, result.data.message

        model = new MapLayerResource(@serialize())
        model.$update success, error

      remove: (cb)->
        return if !confirm('Вы точно хотите удалить этот слой?')
        success =(data)=>
          cb() if cb
          Notify 'success', 'Слой успешно удалён', 3
        error =(result)->
          Message result.data.status, result.data.message

        model = new MapLayerResource({id: @id})
        model.$delete success, error

      addNewImage: ()->
        element = new MapElement(@scope)
        element.inheritPosition()
        element.layer_id = @id
        element.map_id = @map_id
        element.name += " №#{@elements.length + 1}"
        element.prototype_type = 'image'
        element.get_prototype('image')
        FileManager.show({multiple: false}).then (images)=>
          return if images.length == 0
          element.prototype.image = images[0]
          element.create()
          @addToLayer(element)

      addNewText: ()->
        element = new MapElement(@scope)
        element.inheritPosition()
        element.layer_id = @id
        element.map_id = @map_id
        element.name = "New text"
        element.prototype_type = 'text'
        element.get_prototype('text')
        element.prototype.text = "New text"
        element.prototype.style = '{"font-size":"34px"}'
        element.create()
        @addToLayer(element)

      addNewDevice: ()->
      addNewScript: ()->

      addToLayer: (element)->
        @elements.unshift element
        @sortElements()

      sortElements: ()=>
        weight = 0
        for element in @elements
          element.weight = weight
          weight++
        success =(data)->
        error =(result)->
          Message result.data.status, result.data.message
        MapElementResource.sort @elements, success, error

    MapLayer
]