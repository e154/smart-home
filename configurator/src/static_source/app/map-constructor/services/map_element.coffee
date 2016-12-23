angular
.module('angular-map')
.factory 'mapElement', ['$rootScope', '$compile'
  ($rootScope, $compile) ->
    class mapElement

      id: null
      scope: null
      type: null
      name: 'Новый элемент'
      description: ''
      status: 'enabled'
      created_at: null
      update_at: null
      selected: null

      constructor: (@scope)->
        @selected = false

      serialize: ()->
        id: @id if @id
        map_id: map_id if @map_id
        name: @name || ''
        status: @status || ''
        description: @description || ''

      deserialize: (layer)->
        @id = layer?.id || null
        @map_id = layer?.map_id || null
        @name = layer.name || ''
        @description = layer.description || ''
        @status = layer.status || ''

        return @

    mapElement
]