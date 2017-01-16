angular
.module('angular-map')
.factory 'MapScript', ['$rootScope', '$compile', 'Message', 'Notify'
  ($rootScope, $compile, Message, Notify) ->
    class MapScript

      id: null
      scope: null

      constructor: (@scope)->

      serialize: ()->
      deserialize: (data)->
        return @

      create: ()->
        console.log 'create'

      update: ()->
        console.log 'update'

      remove: ()->
        console.log 'remove'

      update: (cb)->
        @upload(cb)

    MapScript
]