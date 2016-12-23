angular
.module('appControllers')
.controller 'mapEditMainWindowCtrl', ['$scope', '$state', 'Message', '$stateParams', 'mapConstructor', 'Notify', '$timeout'
($scope, $state, Message, $stateParams, mapConstructor, Notify, $timeout) ->

  $scope.map.loadEditor('#map-editor1')

  return
]