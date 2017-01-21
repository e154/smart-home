angular
.module('appControllers')
.controller 'deviceNewCtrl', ['$scope', 'Notify', 'Device', '$state', 'Message', 'Node'
($scope, Notify, Device, $state, Message, Node) ->
  vm = this

  vm.devices = []

  Device.group {}, (data)->
    vm.devices = data.devices
    vm.devices.push({name: "Без группы", id: null})
    vm.getNodeInfo()

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
    device: null
    node_id: null
    baud: null
    tty: null
    stop_bite: "2"
    timeout: null
    address: null
    status: "enabled"
  })

  vm.submit =->
    success =(data)->
      $state.go("dashboard.device.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    vm.device.stop_bite = parseInt(vm.device.stop_bite, 10)

    if vm.device?.device?.id != null
      vm.device.stop_bite = null
      vm.device.node_id = null
      vm.device.baud = null
      vm.device.tty = ""
      vm.device.timeout = null
    else
      vm.device.device = null

    vm.device.$create(success, error)

  vm.getNodeInfo =->
    if !vm.device.device?.node?.id
      return

    Node.show {id: vm.device.device.node.id}, (node)->
      vm.device.device.node = node

  vm
]