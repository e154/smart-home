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

  # bpmn scheme editor
  # -----------------------------------------------------------
  scheme4 = bpmnMock.scheme4
  settings =
    engine:
      status: 'editor'
    theme:
      root_path: "/static/themes"

  instance = new bpmnScheme($('#scheme1'))
  instance.setScheme(scheme4)
  instance.setSettings(settings)
  instance.start()

  vm
]