angular
.module('appControllers')
.controller 'deviceNewCtrl', ['$scope', 'Notify', 'Device', '$state', 'Message', 'Node'
($scope, Notify, Device, $state, Message, Node) ->
  vm = this

  Device.get {
    limit:99
    offset: 0
    order: 'desc'
    query: {}
    sortby: 'created_at'
  }, (data)->
    vm.devices = data.devices

  Node.get {
    limit:99
    offset: 0
    order: 'desc'
    query: {}
    sortby: 'created_at'
  }, (data)->
    vm.nodes = data.nodes

  vm.device = new Device({
    name: "Новое устройство"
    description: ""
    device_id: null
    node_id: null
    baud: null
    tty: null
    stop_byte: "2"
    timeout: null
    address: null
    status: "enabled"
  })

  vm.submit =->
    success =(data)->
      $state.go("dashboard.device.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    vm.device.stop_byte = parseInt(vm.device.stop_byte, 10)
    vm.device.$create(success, error)

  vm
]