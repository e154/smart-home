angular
.module('angular-map')
.factory 'MapDeviceState', () ->
  class MapDeviceState

    id: null
    device_state: null
    image: null

    constructor: (@scope, @device_state)->

    serialize: ()->
      return null if !@device_state

      {
        id: @id if @id
        image: {id: @image.id} || null
        device_state: {id: @device_state.id}
      }

    deserialize: (m)->
      @id = m.id if m.id
      @device_state = m.device_state if m.device_state
      @image = m.image if m.image

      @

    remove_image: ()->
      @image = null

  MapDeviceState
