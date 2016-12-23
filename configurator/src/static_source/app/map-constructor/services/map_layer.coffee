angular
.module('angular-map')
.factory 'mapLayer', ['$rootScope', '$compile'
  ($rootScope, $compile) ->
    class mapLayer

      id: null
      map_id: null
      name: 'Новый слой'
      description: ''
      status: null
      elements: null
      scope: null
      created_at: null
      update_at: null

      constructor: (@scope)->

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


    mapLayer
]