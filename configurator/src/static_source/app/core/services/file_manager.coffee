#
# File manager
#

angular
.module('appServices')
.service 'FileManager', [ '$log', 'ngDialog', '$rootScope', 'Stream', '$timeout', 'Notify', '$filter', '$q'
($log, ngDialog, $rootScope, Stream, $timeout, Notify, $filter, $q) ->
  class FileManager
    dialog: null
    scope: null
    title: "File manager"

    constructor: ->
      @scope = $rootScope.$new()
      @scope.file_list = []
      @scope.files_to_upload = []
      @scope.filter_list = []
      @scope.selected_files = []
      @scope.getFileList = @getFileList
      @scope.getFilterList = @getFilterList
      @scope.multiple = false

      @scope.selectFile =(files)=>
        return if !files || !files.length
        angular.forEach files, (file)=>
          @scope.files_to_upload.push file

      $timeout ()=>
        @getFilterList()
      , 2000

    show: (options)=>
      @getFilterList()
      @getFileList()
      defer = $q.defer()
      @scope.defer = defer
      @scope.options = options
      @dialog = ngDialog.open
        template: '/core/templates/file_manager.html'
        controller: 'FileManagerCtrl'
        className: 'ngdialog-theme-default ngdialog-modal-file-manager'
        scope: @scope
        plain: false
        overlay: true
        showClose: false
        closeByDocument: true
        closeByEscape: true
      return defer.promise

    hide: ->
      if @dialog
        @dialog.close()

    getFilterList: =>
      Stream.sendRequest("get_filter_list", {}).then (result)=>
        return if !result.filter_list
        @scope.filter_list = angular.copy(result.filter_list)
        @getFileList(@scope.filter_list[@scope.filter_list.length - 1].date)

    getFileList: (_date)=>
      angular.forEach @scope.filter_list, (f)->
        f.selected = _date == f.date
      Stream.sendRequest("get_image_list", {filter: _date}).then (result)=>
        return if !result.images
        @scope.file_list = angular.copy(result.images)

  new FileManager()
]
