angular
.module('angular-map')
.factory 'mapLayer', ['$rootScope', '$compile'
  ($rootScope, $compile) ->
    class mapLayer

      id: null
      name: "Новый слой"

      elements: null
      scope: null

      constructor: (@scope)->

      serialize: ()->
        {
          id: @id if @id
          name: @name
        }

      deserialize: (layer)->
        @id = layer?.id || null
        @name = layer.name if layer?.name

        return @


    mapLayer
]