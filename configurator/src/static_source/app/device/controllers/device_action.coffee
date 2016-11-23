angular
.module('appControllers')
.controller 'deviceActionCtrl', ['$scope', 'Notify', 'DeviceAction', 'Message', '$stateParams', 'Device'
($scope, Notify, DeviceAction, Message, $stateParams, Device) ->
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
        command: "000000000"
        direction: "inside"
        start_addr: 0
        col_cells: 1
        result_type: "byte"
        function: 2
        description: "Какое-то действие"
        device:
          id: parseInt($stateParams.id, 10)
      })

  vm.submit =->
    success =(result)->
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

  # starting
  # ------------------------------------------
  vm.getDeviceActions()
  vm.getDefaultAction()

  vm
]