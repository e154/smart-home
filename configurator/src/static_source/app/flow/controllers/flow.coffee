angular
.module('appControllers')
.controller 'flowCtrl', ['$scope', 'Notify', 'Flow', 'Stream'
($scope, Notify, Flow, Stream) ->
  vm = this

  $scope.flows = []
  $scope.getStatus = ->
    Stream.sendRequest("get.flows.status", {})

  vm
]