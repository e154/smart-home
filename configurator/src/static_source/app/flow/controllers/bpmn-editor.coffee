angular
.module('appControllers')
.controller 'bpmnEditorCtrl', ['$scope', 'Notify', 'Flow', '$stateParams', '$state', '$timeout', 'bpmnMock'
'bpmnScheme', 'bpmnSettings', 'log'
($scope, Notify, Flow, $stateParams, $state, $timeout, bpmnMock, bpmnScheme, bpmnSettings, log) ->
  vm = this

  # settings
  #------------------------------------------------------------------------------
  settings =
    engine:
      status: 'editor'
    theme:
      root_path: "/static/themes"
    editorPallet:
      groups: [
        {
          name: 'event'
          items: [
            {
              type:
                name: 'event'
                start:
                  0:
                    0: true
              title: 'start'
              class: 'bpmn-icon-start-event-none'
              tooltip: 'Create start event'
              shape: bpmnSettings.template('event')
            }
            {
              type:
                name: 'event'
                intermediate:
                  0:
                    0: true
              title: 'intermediate'
              class: 'bpmn-icon-intermediate-event-none'
              tooltip: 'Create intermediate event'
              shape: bpmnSettings.template('event')
            }
            {
              type:
                name: 'event'
                end:
                  simply:
                    top_level: true
              title: 'end'
              class: 'bpmn-icon-end-event-none'
              tooltip: 'Create end event'
              shape: bpmnSettings.template('event')
            }
          ]
        }
        {
          name: 'gateway'
          items: [
            {
              type:
                name: 'gateway'
                start:
                  0:
                    0: true
              title: 'gateway'
              class: 'bpmn-icon-gateway-xor'
              tooltip: 'Create gateway'
              shape: bpmnSettings.template('gateway')
            }
          ]
        }
        {
          name: 'task'
          items: [
            {
              type:
                name: 'task'
              title: 'task'
              class: 'bpmn-icon-task-none'
              tooltip: 'Create task'
              shape: bpmnSettings.template('task')
            }
          ]
        }
        {
          name: 'group'
          items: [
            {
              type:
                name: 'group'
              title: 'group'
              class: 'bpmn-icon-subprocess-expanded'
              tooltip: 'Create group'
              shape: bpmnSettings.template('group')
            }
          ]
        }
      ]

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

  $scope.serialise =->
    $scope.scheme = redactor.getScheme()

  $scope.callback['save']= ()->
    $scope.serialise()

  $timeout ()->
    $scope.$apply(
      $scope.callback
    )

  vm
]