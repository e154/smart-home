angular
.module('angular-map')
.factory 'MapDevice', ['$rootScope', '$compile', 'Message', 'Notify'
  ($rootScope, $compile, Message, Notify) ->
    class MapDevice

      id: null
      scope: null
      device: null

      constructor: (@scope)->

      serialize: ()->
        id: @id if @id

      deserialize: (m)->
        @id = m.id if m.id

        @

      create: ()->
      update: ()->
      remove: ()->
      update: (cb)->
        @upload(cb)

    MapDevice
]