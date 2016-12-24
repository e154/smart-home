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

      constructor: (@scope, @layer_id)->

      serialize: ()->
        id: @id
        map_id: @map_id
        layer_id: @layer_id
        name: @name
        status: @status
        description: @description || ''
        created_at: @created_at if @created_at
        update_at: @update_at if @update_at

      deserialize: (element)->
        @id = element.id || null
        @map_id = element.map.id || null
        @layer_id = element.layer.id || null
        @name = element.name || ''
        @description = element.description || ''
        @status = element.status || ''
        @type = element.type || 'image'

        return @

    mapElement
]