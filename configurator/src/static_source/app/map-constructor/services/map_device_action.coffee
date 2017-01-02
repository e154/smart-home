angular
.module('angular-map')
.factory 'MapDeviceAction', () ->
  class MapDeviceAction

    id: null

    constructor: (@scope)->

    serialize: ()->

      {
        id: @id if @id
      }

    deserialize: (m)->
      @id = m.id if m.id

      @

  MapDeviceAction
