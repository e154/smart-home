angular
.module('appControllers')
.controller 'bpmnEditorCtrl', ['$scope', 'Notify', 'Flow', '$stateParams', '$state', '$timeout', 'bpmnMock'
'bpmnScheme', 'bpmnSettings', '$http', 'log', 'Worker', 'ngDialog', '$filter', 'Stream'
($scope, Notify, Flow, $stateParams, $state, $timeout, bpmnMock, bpmnScheme, bpmnSettings, $http
log, Worker, ngDialog, $filter, Stream) ->
  vm = this

  $scope.selected = []
  $scope.selectedConn =
    title: ""
    object: ""
    direction: ""
  $scope.directions = [
    {
      name: $filter('translate')('true')
      value: "true"
    }
    {
      name: $filter('translate')('false')
      value: "false"
    }
    {
      name: $filter('translate')('no mater')
      value: ""
    }
  ]

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
        {
          name: 'flow'
          items: [
            {
              type:
                name: 'flow'
              title: 'flow'
              class: 'bpmn-icon-participant'
              tooltip: 'Create flow link'
              shape: bpmnSettings.template('flow')
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
  redactor.scope.$watch 'selected', (selected)=>
    return if !selected

    $scope.selected = []
    connections = redactor.getAllConnections()
    connection = null
    for object in selected
      angular.forEach $scope.redactor.scope.intScheme.objects, (obj, key)->
        if key == object.id
          $scope.selected.push obj
          if !$scope.elementScripts.hasOwnProperty(obj.data.id)
            $scope.elementScripts[obj.data.id] = null

          if !$scope.elementFlows.hasOwnProperty(obj.data.id)
            $scope.elementFlows[obj.data.id] = null

      # connections
      if object.type != "connector"
        continue

      angular.forEach connections, (conn)->
        if object.id == conn.id
          connection = conn

    if connection?
      $scope.selectedConn =
        title: redactor.getLabel(connection)
        direction: connection.getParameter("direction") || ""
        object: connection
    else $scope.selectedConn =
        object: null

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

  $scope.doWorker =(worker, $index)->
    return if !worker.id
    Stream.sendRequest("do.worker", {worker_id: worker.id}).then (result)->
      if !result.error
        Notify 'success', "Команда выполнена успешно", 3
      else
        Notify 'error', "Результат выполнения команды:\n\r #{result.error}", 3

  # get device actions (select2)
  $scope.actions = []
  $scope.refreshActions = (query)->
    $http(
      method: 'GET'
      url: window.app_settings.server_url + "/api/v1/device_action/search"
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
      url: window.app_settings.server_url + "/api/v1/script/search"
      params:
        query: query
        limit: 5
        offset: 0
    ).then (response)->
      $scope.scripts = response.data.scripts

  # scripts
  #------------------------------------------------------------------------------
  $scope.showScript =(script, e)->
    e.preventDefault()
    $scope.script = script

    ngDialog.open
      scope: $scope
      showClose: false
      template: '/script/templates/modal.show.html'
      className: 'ngdialog-theme-default ngdialog-scripts-show'
      controller: 'scriptModalShowCtrl'
      controllerAs: 'script'

  $scope.addScript =(id, e)->
    e.preventDefault()
    $scope.setScript = (script)->
      $scope.elementScripts[id] = script

    ngDialog.open
      scope: $scope
      showClose: false
      closeByEscape: false
      closeByDocument: false
      template: '/script/templates/modal.new.html'
      className: 'ngdialog-theme-default ngdialog-scripts-edit'
      controller: 'scriptModalNewCtrl'
      controllerAs: 'script'

  $scope.editScript =(elementScripts, id, e)->
    e.preventDefault()
    $scope.script = elementScripts[id]
    $scope.setScript = (script)->
      $scope.elementScripts[id] = script

    ngDialog.open
      scope: $scope
      showClose: false
      closeByEscape: false
      closeByDocument: false
      template: '/script/templates/modal.edit.html'
      className: 'ngdialog-theme-default ngdialog-scripts-edit'
      controller: 'scriptModalEditCtrl'
      controllerAs: 'script'

  # connections
  #------------------------------------------------------------------------------
  $scope.setDirection =->
    return if !$scope.selectedConn
    direction = $scope.selectedConn.direction
    $scope.selectedConn.object.setParameter("direction", direction)

  $scope.setLabel =->
    return if !$scope.selectedConn
    conn = $scope.selectedConn.object
    redactor.setLabel(conn, $scope.selectedConn.title)

  # flows
  #------------------------------------------------------------------------------
  # select flows for flow elements (select2)
  $scope.flows = []
  $scope.refreshFlows = (query)->
    $http(
      method: 'GET'
      url: window.app_settings.server_url + "/api/v1/flow/search"
      params:
        query: query
        limit: 5
        offset: 0
    ).then (response)->
      $scope.flows = []
      angular.forEach response.data.flows, (flow)->
        if flow.id != $scope.flow.id
          $scope.flows.push flow

  vm
]