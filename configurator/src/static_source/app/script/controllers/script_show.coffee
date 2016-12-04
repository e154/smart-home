angular
.module('appControllers')
.controller 'scriptShowCtrl', ['$scope', 'Notify', 'Script', '$stateParams', '$state', '$timeout'
($scope, Notify, Script, $stateParams, $state, $timeout) ->
  vm = this

  success = (script) ->
    vm.script = script
    $scope.ace_options.readOnly = true

  error = ->
    $state.go 'dashboard.script.index'

  Script.show {id: $stateParams.id}, success, error

  $scope.$watch 'script.script.lang', (lang)->
    return if !lang || lang == ''
    switch lang
      when 'javascript'
        $scope.ace_options.mode = 'javascript'
      when 'coffeescript'
        $scope.ace_options.mode = 'coffee'
      when 'lua'
        $scope.ace_options.mode = 'lua'

  vm
]