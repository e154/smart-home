angular
.module('angular-map')
.factory 'mapElement', ['$rootScope', '$compile'
  ($rootScope, $compile) ->
    class mapElement

      scope: null

      id: null
      layer_id: null
      map_id: null
      name: 'Новый элемент'
      description: ''
      type: 'image'
      status: 'enabled'
      selected: false
      created_at: null
      update_at: null
      weight: 0

      constructor: (@scope, @layer_id)->

      serialize: ()->
        name: @name
        id: @id if @id
        map: {id: @map_id} if @map_id
        layer: {id: @layer_id} if @layer_id
        status: @status
        description: @description
        created_at: @created_at if @created_at
        update_at: @update_at if @update_at
        weight: @weight
        type: @type

      deserialize: (element)->
        @id = element.id || null
        @map_id = element.map.id || null
        @layer_id = element.layer.id || null
        @name = element.name || ''
        @description = element.description || ''
        @status = element.status || ''
        @type = element.type || 'image'
        @weight = element.weight || 0
        @created_at = element.created_at || ''
        @update_at = element.update_at || ''

        return @

    mapElement
]