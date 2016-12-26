angular
.module('angular-map')
.factory 'MapDevice', ['$rootScope', '$compile', 'Message', 'Notify'
  ($rootScope, $compile, Message, Notify) ->
    class MapDevice

      id: null
      scope: null

      constructor: (@scope)->

      serialize: ()->
      deserialize: (data)->
        return @

      create: ()->
      update: ()->
      remove: ()->
      update: (cb)->
        @upload(cb)

    MapDevice
]