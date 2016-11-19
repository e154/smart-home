angular
.module('appControllers')
.controller 'bpmnViewerCtrl', ['$scope', 'Notify', 'Flow', '$stateParams', '$state', '$timeout', 'bpmnMock'
'bpmnScheme', 'bpmnSettings'
($scope, Notify, Flow, $stateParams, $state, $timeout, bpmnMock, bpmnScheme, bpmnSettings) ->
  vm = this

  # settings
  #------------------------------------------------------------------------------
  settings =
    engine:
      container:
        zoom: true
    theme:
      root_path: "/static/themes"

  # redactor
  #------------------------------------------------------------------------------
  redactor = new bpmnScheme($('#scheme1'))
  redactor.setSettings(settings)
  redactor.start()

  # prepare scheme
  #------------------------------------------------------------------------------
  $scope.$watch 'flow', (scheme)->
    if !scheme || !scheme?.name
      return
    redactor.setScheme(scheme)
    redactor.restart()

  vm
]