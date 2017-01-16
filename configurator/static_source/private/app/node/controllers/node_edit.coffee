angular
.module('appControllers')
.controller 'nodeEditCtrl', ['$scope', 'Message', '$stateParams', 'Node', '$state'
($scope, Message, $stateParams, Node, $state) ->
  vm = this

  Node.show {id: $stateParams.id}, (node)->
    vm.node = node

  vm.remove =->
    if confirm('точно удалить узел?')
      remove()

  remove =->
    success =->
      $state.go("dashboard.node.index")
    error =(result)->
      Message result.data.status, result.data.message
    vm.node.$delete success, error

  vm.submit =->
    success =(data)->
      $state.go("dashboard.node.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    vm.node.$update(success, error)

  vm
]