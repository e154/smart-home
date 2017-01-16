angular
.module('appControllers')
.controller 'scriptCtrl', ['$scope', 'Notify', 'Script', 'Stream', '$log'
($scope, Notify, Script, Stream, $log) ->
  vm = this

  $scope.ace_options =
    useWrapMode: true
    mode:'coffee'
    theme:'dawn'
    advanced:{}
    workerPath:'/static/js/ace-builds/src-noconflict'

  vm
]