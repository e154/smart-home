angular
.module('angular-map')
.factory 'MapDeviceAction', () ->
  class MapDeviceAction

    id: null
    image: null
    type: null

    constructor: (@scope, @device_action)->

    serialize: ()->

      {
        id: @id if @id
        image: @image
        device_action: {id: @device_action.id}
        type: @type || ''
      }

    removeImage: ()->
      @image = null

    deserialize: (m)->
      @id = m.id if m.id
      @image = m.image if m.image
      @type = m.type || ''

      @

  MapDeviceAction
