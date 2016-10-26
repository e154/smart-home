angular
.module('appControllers')
.controller 'nodeEditCtrl', ['$scope', 'Notify', '$stateParams', 'Node', '$state'
($scope, Notify, $stateParams, Node, $state) ->
  vm = this

  Node.show {id: $stateParams.id}, (node)->
    vm.node = node

  vm.submit =->
    success =(data)->
      $state.go("dashboard.node.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    vm.node.$update(success, error)

  vm
]