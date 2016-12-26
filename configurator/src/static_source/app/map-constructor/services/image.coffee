angular
.module('angular-map')
.factory 'Image', ['$rootScope', '$compile', 'Message', 'Notify', 'Upload', '$timeout'
  ($rootScope, $compile, Message, Notify, Upload, $timeout) ->
    class Image

      file: null
      progress: null

      constructor: (@file)->

      upload: (cb)=>
        upload = Upload.upload
          url: window.server_url + "/api/v1/image/upload"
          data:
            files: [@file]

        success =(response)->
          cb(response.data) if cb
          Notify 'success', 'Файл загружен', 3

        error =(response)->
          if response.status > 0
            Notify 'error', response.status + ': ' + response.data

        progress =(evt)=>
          @progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total))

        upload.then success, error, progress

    Image
]