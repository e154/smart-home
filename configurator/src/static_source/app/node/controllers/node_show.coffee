angular
.module('appControllers')
.controller 'nodeShowCtrl', ['$scope', 'Notify', 'Node', '$stateParams'
($scope, Notify, Node, $stateParams) ->
  vm = this

  vm.node = Node.get {id: $stateParams.id}

  vm
]