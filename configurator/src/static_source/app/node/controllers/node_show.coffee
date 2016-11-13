angular
.module('appControllers')
.controller 'nodeShowCtrl', ['$scope', 'Notify', 'Node', '$stateParams', '$state'
($scope, Notify, Node, $stateParams, $state) ->
  vm = this

  success = (node) ->
    vm.node = node

  error = ->
    $state.go 'dashboard.node.index'

  Node.show {id: $stateParams.id}, success, error


  vm
]