angular
.module('appControllers')
.controller 'scriptShowCtrl', ['$scope', 'Notify', 'Script', '$stateParams', '$state', 'Message'
($scope, Notify, Script, $stateParams, $state, Message) ->
  vm = this

  success = (script) ->
    vm.script = script
    $scope.ace_options.readOnly = true

  error = ->
    Message result.data.status, result.data.message

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