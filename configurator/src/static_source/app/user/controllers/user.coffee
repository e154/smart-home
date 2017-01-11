angular
.module('appControllers')
.controller 'userCtrl', ['$scope', '$translate'
($scope, $translate) ->

  $translate.refresh()
]