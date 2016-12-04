angular
.module('appControllers')
.controller 'deviceActionCtrl', ['$scope', 'Notify', 'DeviceAction', 'Message', '$stateParams', 'Device'
'$http'
($scope, Notify, DeviceAction, Message, $stateParams, Device, $http) ->
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

    vm.current.$remove(success, error)

  # select2
  # ------------------
  $scope.scripts = []
  $scope.refreshScripts = (query)->
    $http(
      method: 'GET'
      url: window.server_url + "/api/v1/script/search"
      params:
        query: query
        limit: 5
        offset: 0
    ).then (response)->
      $scope.scripts = response.data.scripts

  # starting
  # ------------------------------------------
  vm.getDeviceActions()
  vm.getDefaultAction()

  vm
]