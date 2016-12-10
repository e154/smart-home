angular
.module('appControllers')
.controller 'bpmnEditorCtrl', ['$scope', 'Notify', 'Flow', '$stateParams', '$state', '$timeout', 'bpmnMock'
'bpmnScheme', 'bpmnSettings', '$http', 'log'
($scope, Notify, Flow, $stateParams, $state, $timeout, bpmnMock, bpmnScheme, bpmnSettings, $http, log) ->
  vm = this

  $scope.selected = []

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
      ]


  # redactor
  #------------------------------------------------------------------------------
  $scope.redactor = redactor = new bpmnScheme($('#scheme1'))
  redactor.setSettings(settings)

  $scope.$watch 'flow', (scheme)->
    if !scheme || !scheme?.name
      return
    redactor.start()
    redactor.setScheme(scheme)

  $scope.serialise =->
    $scope.scheme = redactor.getScheme()

  $scope.callback['save']= ()->
    $scope.serialise()

  $timeout ()->
    $scope.$apply(
      $scope.callback
    )

  $scope.removeElement =(element, $index)->
    index = $scope.selected.indexOf(element)
    object =
      id: element.data.id
      type: 'shape'
    $scope.redactor.removeObject(object)

    if index > -1
      $scope.selected.splice(index, 1)

  # select element on scheme
  #------------------------------------------------------------------------------
  redactor.scope.$watch 'selected', (objects)=>
    return if !objects

    $scope.selected = []
    for object in objects
      angular.forEach $scope.redactor.scope.intScheme.objects, (obj, key)->
        if key == object.id
          $scope.selected.push obj
          if !$scope.elementScripts.hasOwnProperty(obj.data.id)
            $scope.elementScripts[obj.data.id] = null

    $timeout ()->
      $scope.$apply()

  , true

  # workers
  #------------------------------------------------------------------------------
  $scope.addNewWorker =->
    worker = new Worker({
      name: 'Действие'
      time: '* * * * * *'
      status: 'enabled'
      flow:
        id: parseInt($stateParams.id, 10)
      device_action: null
      workflow:
        id: $scope.flow.workflow.id
    })

    $scope.flow.workers.push worker

  $scope.removeWorker =(worker, $index)->
    if !worker.id
      $scope.flow.workers.splice($index, 1)
      return

    for i in [0...$scope.flow.workers.length]
      if $scope.flow.workers[i].id == worker.id
        $scope.flow.workers.splice(i, 1)
        break

  # get device actions (select2)
  $scope.actions = []
  $scope.refreshActions = (query)->
    $http(
      method: 'GET'
      url: window.server_url + "/api/v1/device_action/search"
      params:
        query: query
        limit: 5
        offset: 0
    ).then (response)->
      $scope.actions = response.data.actions

  # select scripts for flow elements (select2)
  $scope.scripts = []
  $scope.refreshScripts = (query)->
    $http(
      method: 'GET'
      url: window.server_url + "/api/v1/script/search"
      params:
        query: query
        limit: 5
        offset: 0
    ).then (response)->
      $scope.scripts = response.data.scripts

  vm
]