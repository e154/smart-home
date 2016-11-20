angular
.module('appControllers')
.controller 'deviceEditCtrl', ['$scope', 'Message', '$stateParams', 'Device', '$state', 'Node'
($scope, Message, $stateParams, Device, $state, Node) ->
  vm = this

  vm.nodes = {}
  vm.devices = []

  getDevices =->
    Device.group {}, (data)->
      angular.forEach data.devices, (device)->
        if device.id != vm.device.id
          vm.devices.push(device)
      vm.devices.push({name: "Без группы", id: null})

  Node.get {
    limit:99
    offset: 0
    order: 'desc'
    query: {}
    sortby: 'created_at'
  }, (data)->
    vm.nodes = data.nodes
    getDevices()

  Device.show {id: $stateParams.id}, (device)->
    vm.device = device
    vm.device.stop_bite = device.stop_bite.toString()
    vm.getNodeInfo()

  vm.remove =->
    if confirm('точно удалить узел?')
      remove()

  remove =->
    success =->
      $state.go("dashboard.device.index")
    error =(result)->
      Message result.data.status, result.data.message
    vm.device.$delete success, error

  vm.submit =->
    success =(data)->
      $state.go("dashboard.device.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    vm.device.stop_bite = parseInt(vm.device.stop_bite, 10)

    if vm.device.device?
      vm.device.stop_bite = null
      vm.device.node_id = null
      vm.device.baud = null
      vm.device.tty = ""
      vm.device.timeout = null
    else
      vm.device.device = null

    vm.device.$update(success, error)

  vm.getNodeInfo =->
    if !vm.device.device?.node?.id
      return

    Node.show {id: vm.device.device.node.id}, (node)->
      vm.device.device.node = node

  vm
]