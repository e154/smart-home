angular
.module('appControllers')
.controller 'nodeNewCtrl', ['$scope', 'Notify', 'Node', '$state', 'Message'
($scope, Notify, Node, $state, Message) ->
  vm = this

  vm.node = new Node({
    name: "node"
    ip: "127.0.0.1"
    port: 3000
    description: ""
  })

  vm.submit =->
    success =(data)->
      $state.go("dashboard.node_show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    vm.node.$create(success, error)

  vm
]