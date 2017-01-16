angular
.module('appControllers')
.controller 'scriptModalShowCtrl', ['$scope', 'Notify', 'Script', '$stateParams', '$state', 'Message'
($scope, Notify, Script, $stateParams, $state, Message) ->
  vm = this
  vm.script = $scope.$parent.script

  $scope.ace_options =
    useWrapMode: true
    mode:'coffee'
    theme:'dawn'
    advanced:{}
    workerPath:'/static/js/ace-builds/src-noconflict'
    readOnly: true

  switch vm.script.lang
    when 'javascript'
      $scope.ace_options.mode = 'javascript'
    when 'coffeescript'
      $scope.ace_options.mode = 'coffee'
    when 'lua'
      $scope.ace_options.mode = 'lua'

  vm
]