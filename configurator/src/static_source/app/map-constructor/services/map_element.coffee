angular
.module('angular-map')
.factory 'MapElement', ['$rootScope', '$compile', 'MapElementResource', 'Message'
'Notify', 'MapImage', 'MapDevice', 'MapScript', 'MapText'
  ($rootScope, $compile, MapElementResource, Message, Notify, MapImage, MapDevice, MapScript, MapText) ->
    class MapElement

      scope: null

      prototype: null
      old_prototype: null

      current_tab: 'base'
      id: null
      layer_id: null
      map_id: null
      name: 'Новый элемент'
      description: ''
      prototype_type: ''
      prototype_id: null
      prototype: null
      status: 'enabled'
      selected: false
      created_at: null
      update_at: null
      weight: 0
      graph_settings:
        width: null
        height: null
        position:
          top: 0
          left: 0

      constructor: (@scope, @layer_id)->
        @scope.$watch(()=>
            @prototype_type
        , (val, old_val)=>
          return if val == old_val
          @get_prototype()
        )
        @get_prototype()

      resetPosition: ()->
        @graph_settings.position =
          top: 0
          left: 0

      inheritPosition: ()->
        position = $(".map-editor").position()
        @graph_settings.position =
          top: parseInt(position.top, 10) * -1
          left: parseInt(position.left, 10) * -1

      serialize: ()->
        prototype = @prototype?.serialize() || null
        name: @name
        id: @id if @id
        map: {id: @map_id} if @map_id
        layer: {id: @layer_id} if @layer_id
        status: @status
        description: @description
        created_at: @created_at if @created_at
        update_at: @update_at if @update_at
        weight: @weight
        prototype_type: @prototype_type
        prototype_id: @prototype_id if @prototype_id
        graph_settings: angular.toJson(@graph_settings)
        prototype: prototype if prototype

      deserialize: (element)->
        @id = element.id || null
        @map_id = element.map.id if element.map?.id
        @layer_id = element.layer.id if element.layer?.id
        @name = element.name || ''
        @description = element.description || ''
        @status = element.status || ''
        @prototype_type = element.prototype_type || ''
        @prototype_id = element.prototype_id if element.prototype_id
        @weight = element.weight || 0
        @created_at = element.created_at || ''
        @update_at = element.update_at || ''
        @graph_settings = angular.fromJson(element.graph_settings) if element.graph_settings

        if element.prototype
          @get_prototype(element.prototype)

        return @

      create: (cb)->
        success =(data)=>
          @id = data.id
          Notify 'success', 'Элемент успешно создан', 3
          cb() if cb
        error =(result)->
          Message result.data.status, result.data.message

        model = new MapElementResource(@serialize())
        model.$create success, error

      update_element_only: (cb)->
        update: (cb)->
        success =(data)=>
          cb() if cb
        error =(result)->
          Message result.data.status, result.data.message

        model = new MapElementResource(@serialize())
        model.$update_element_only success, error

      update: (cb)->
        success =(data)=>
          Notify 'success', 'Элемент успешно обновлён', 3
          cb() if cb
        error =(result)->
          Message result.data.status, result.data.message

        model = new MapElementResource(@serialize())
        model.$update success, error

      remove: (cb)->
        return if !confirm('Вы точно хотите удалить этот элемент?')
        success =(data)=>
          cb() if cb
          Notify 'success', 'Элемент успешно удалён', 3
        error =(result)->
          Message result.data.status, result.data.message

        model = new MapElementResource({id: @id})
        model.$delete success, error

      get_prototype: (data)->
        switch @prototype_type
          when 'text'
            @prototype = new MapText(@scope)
          when 'image'
            @prototype = new MapImage(@scope)
          when 'device'
            @prototype = new MapDevice(@scope)
          when 'script'
            @prototype = new MapScript(@scope)

        if data
          @prototype.deserialize(data)

      copy: (_element)->
        @deserialize(_element.serialize())
        @id = null
        @prototype_id = null
        @name = "#{@name} (copy)" if @name.indexOf('(copy)') == -1

        @

    MapElement
]