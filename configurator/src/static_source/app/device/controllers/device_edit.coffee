angular
.module('appControllers')
.controller 'deviceEditCtrl', ['$scope', 'Message', '$stateParams', 'Device', '$state', 'Node'
($scope, Message, $stateParams, Device, $state, Node) ->
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

  Device.show {id: $stateParams.id}, (device)->
    vm.device = device

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

    vm.device.$update(success, error)

  vm
]