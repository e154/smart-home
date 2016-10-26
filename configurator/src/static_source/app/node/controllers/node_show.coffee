angular
.module('appControllers')
.controller 'nodeShowCtrl', ['$scope', 'Notify', 'Node', '$stateParams'
($scope, Notify, Node, $stateParams) ->
  vm = this

  Node.show {id: $stateParams.id}, (node)->
    vm.node = node


  vm
]