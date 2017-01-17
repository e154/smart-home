angular
.module('appControllers')
.controller 'deviceActionCtrl', ['$scope', 'Notify', 'DeviceAction', 'Message', '$stateParams', 'Device'
'$http', 'ngDialog', 'Stream'
($scope, Notify, DeviceAction, Message, $stateParams, Device, $http, ngDialog, Stream) ->
  vm = this
  vm.actions = []
  vm.current ={}
  vm.last_current =null

  vm.addNew =->
    vm.getDefaultAction()

  vm.show =(action)->
    vm.last_current = angular.copy action
    vm.current = new DeviceAction(action)

  vm.getDeviceActions =->
    Device.actions {id: $stateParams.id}, (actions)->
      vm.actions = actions

  vm.getDefaultAction =->
    vm.current = new DeviceAction({
        name: "Новое действие"
        command: "03000005"
        script: null
        description: "Какое-то действие"
        device:
          id: parseInt($stateParams.id, 10)
      })

  vm.submit =->
    success =(result)->
      Notify 'success', 'Действие успешно обновлено', 3

      vm.getDeviceActions()
      vm.getDefaultAction()

    error =(result)->
      Message result.data.status, result.data.message

    if !vm.current.id
      vm.current.$create(success, error)
    else
      vm.current.$update(success, error)

  vm.cancel =->
    return if !vm.last_current
    vm.current = new DeviceAction(vm.last_current)

  vm.remove =->
    success =->
      vm.getDeviceActions()
      vm.getDefaultAction()

    error =(result)->
      Message result.data.status, result.data.message

    if confirm('Вы точно хотите удалить это действие?')
      vm.current.$remove(success, error)

  # select2
  # ------------------
  $scope.scripts = []
  $scope.refreshScripts = (query)->
    $http(
      method: 'GET'
      url: window.app_settings.server_url + "/api/v1/script/search"
      params:
        query: query
        limit: 5
        offset: 0
    ).then (response)->
      $scope.scripts = response.data.scripts

  # scripts
  #------------------------------------------------------------------------------
  vm.showScript =(script, e)->
    e.preventDefault()
    $scope.script = script

    ngDialog.open
      scope: $scope
      showClose: false
      template: '/script/templates/modal.show.html'
      className: 'ngdialog-theme-default ngdialog-scripts-show'
      controller: 'scriptModalShowCtrl'
      controllerAs: 'script'

  vm.addScript =(e)=>
    e.preventDefault()
    $scope.setScript =(script)=>
      vm.current.script = script

    ngDialog.open
      scope: $scope
      showClose: false
      closeByEscape: false
      closeByDocument: false
      template: '/script/templates/modal.new.html'
      className: 'ngdialog-theme-default ngdialog-scripts-edit'
      controller: 'scriptModalNewCtrl'
      controllerAs: 'script'

  vm.editScript =(script, e)->
    e.preventDefault()
    $scope.script = script
    $scope.setScript =(script)=>
      vm.current.script = script

    ngDialog.open
      scope: $scope
      showClose: false
      closeByEscape: false
      closeByDocument: false
      template: '/script/templates/modal.edit.html'
      className: 'ngdialog-theme-default ngdialog-scripts-edit'
      controller: 'scriptModalEditCtrl'
      controllerAs: 'script'

  vm.doAction =(action, e)->
    e.preventDefault()
    e.stopPropagation()
    return if !action.id

    Stream.sendRequest("do.action", {action_id: action.id, device_id: parseInt($stateParams.id, 10)}).then (result)->
      if !result.error
        Notify 'success', "Команда выполнена успешно", 3
      else
        Notify 'error', "Результат выполнения команды:\n\r #{result.error}", 3

  # starting
  # ------------------------------------------
  vm.getDeviceActions()
  vm.getDefaultAction()

  vm
]