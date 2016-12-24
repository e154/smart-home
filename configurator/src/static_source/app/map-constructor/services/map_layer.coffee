angular
.module('angular-map')
.factory 'mapLayer', ['$rootScope', '$compile', 'mapElement'
  ($rootScope, $compile, mapElement) ->
    class mapLayer

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
          @elements.push new mapElement(@scope).deserialize(element)

        return @

      addElement: (element)=>
        @elements.push element

    mapLayer
]