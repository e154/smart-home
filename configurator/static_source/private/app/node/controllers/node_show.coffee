angular
.module('appControllers')
.controller 'nodeShowCtrl', ['$scope', 'Notify', 'Node', '$stateParams', '$state', '$timeout'
($scope, Notify, Node, $stateParams, $state, $timeout) ->
  vm = this

  success = (node) ->
    vm.node = node
    $timeout ()->
      $scope.getStatus().then (result)->
        $scope.nodes = result.nodes

        angular.forEach $scope.nodes, (value, id)->
          if node.id == parseInt(id, 10)
            vm.node.state = value
    , 500

  error = ->
    $state.go 'dashboard.node.index'

  Node.show {id: $stateParams.id}, success, error

  vm
]