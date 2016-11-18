angular
.module('appControllers')
.controller 'nodeCtrl', ['$scope', 'Notify', 'Node', 'Stream', '$log'
($scope, Notify, Node, Stream, $log) ->
  vm = this

  $scope.nodes = []
  $scope.getStatus = ->
    Stream.sendRequest("get.nodes.status", {})

  vm
]