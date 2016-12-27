angular
.module('angular-map')
.factory 'MapImage', ['$rootScope', '$compile', 'Message', 'Notify', 'Image'
  ($rootScope, $compile, Message, Notify, Image) ->
    class MapImage

      id: null
      scope: null
      style: ''
      image: null
      file: null

      constructor: (@scope)->

      remove_image: ()->
        @image = null

      serialize: ()->
        return null if !@image
        id: @id if @id
        image:
          id: @image.id
        style: @style

      deserialize: (m)->
        @id = m.id if m.id
        @image = m.image if m.image
        @style = m.style if m.style

        @

    MapImage
]