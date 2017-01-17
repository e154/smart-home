###*
# Created by delta54 on 20.12.14.
###

angular
.module('appControllers')
.controller 'FileManagerCtrl', ['$scope','Notify','$log','Stream','$timeout','Upload'
  ($scope, Notify, $log, Stream, $timeout, Upload) ->
    vm = this

    $scope.mode = "select" # select|upload

    $scope.select = ($event, file)=>
      $event.preventDefault()
      index = $scope.selected_files.indexOf(file)
      if index == -1
        if !$scope.options.multiple && $scope.selected_files.length
          $scope.selected_files[0].selected = false
          $scope.selected_files.splice(0, 1)

        $scope.selected_files.push(file)
        file.selected = true
      else
        $scope.selected_files.splice(index, 1)
        file.selected = false

    $scope.removeFile = ($event, f)=>
      $event.stopPropagation()
      $event.preventDefault()
      return if !f

      if $scope.mode == "select"
        Stream.sendRequest("remove_image", {image_id: f.id}).then (result)=>
          if result.status == 'ok'
            index = $scope.file_list.indexOf(f)
            if index > -1
              $scope.file_list.splice(index, 1)
          if $scope.file_list.length == 0
            $scope.getFilterList()

      else if $scope.mode == "upload"
        angular.forEach $scope.files_to_upload, (file, key)=>
          if file.$$hashKey == f.$$hashKey
            $scope.files_to_upload.splice(key, 1)

        if $scope.files_to_upload.length == 0
          $scope.getFilterList()

    $scope.submit =->
      files = []
      angular.forEach $scope.selected_files, (file)->
        files.push file
      $scope.defer.resolve(files)


    $scope.upload =->

      if !$scope.files_to_upload || !$scope.files_to_upload.length
        return

      upload = Upload.upload(
        url: window.app_settings.server_url + "/api/v1/image/upload"
        data: files: $scope.files_to_upload
      )

      success =->
        $timeout ->
          clear()
          $scope.getFilterList()
          $scope.mode = "select"

      error =(response) ->
        if response.status > 0
          Notify 'error', response.status + ': ' + response.data

      progress =(evt) ->
        progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total))

      upload.then success, error, progress

    clear =->
      $scope.files_to_upload = []
      $scope.selected_files = []

    return
]
