angular
.module('angular-map')
.factory 'mapElement', ['$rootScope', '$compile'
  ($rootScope, $compile) ->
    class mapElement

      id: null
      layer_id: null
      map_id: null
      scope: null
      type: 'image'
      name: 'Новый элемент'
      description: ''
      status: 'enabled'
      created_at: null
      update_at: null
      selected: false

      constructor: (@scope, @layer_id)->

      serialize: ()->
        id: @id if @id
        map_id: map_id if @map_id
        layer_id: @layer_id
        name: @name || ''
        status: @status || ''
        description: @description || ''

      deserialize: (element)->
        @id = element?.id || null
        @map_id = element?.map_id || null
        @layer_id = element?.layer_id || null
        @name = element.name || ''
        @description = element.description || ''
        @status = element.status || ''
        @type = element.type || 'image'

        return @

    mapElement
]