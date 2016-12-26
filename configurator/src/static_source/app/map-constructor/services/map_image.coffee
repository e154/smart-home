angular
.module('angular-map')
.factory 'MapImage', ['$rootScope', '$compile', 'Message', 'Notify', 'Image', 'FileManager'
  ($rootScope, $compile, Message, Notify, Image, FileManager) ->
    class MapImage

      id: null
      scope: null
      style: ''
      image: null
      file: null

      constructor: (@scope)->

      remove_image: ()->
        @image = null

      show_file_manager: ()->
        console.log 'show_file_manager'
        FileManager.show().then (image)=>
          @image = image[0]

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