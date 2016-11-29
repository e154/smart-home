angular
.module('appControllers')
.controller 'flowEditCtrl', ['$scope', 'Message', '$stateParams', 'Flow', '$state', 'Workflow', '$timeout'
'log', 'Notify', 'Worker', '$http'
($scope, Message, $stateParams, Flow, $state, Workflow, $timeout, log, Notify, Worker
$http) ->
  vm = this

  # vars
  $scope.callback = {}
  $scope.workflows = []
  $scope.flow = {}

  # workflow list
  #------------------------------------------------------------------------------
  getWorkflow =->
    success = (result)->
      $scope.workflows = result.items
    error = (result)->
      Message result.data.status, result.data.message
    Workflow.all {}, success, error

  # get flow
  #------------------------------------------------------------------------------
  getFlow =->
    success = (flow) ->
      $scope.flow = flow
      if !$scope.flow?.workers
        $scope.flow.workers = []

      $timeout ()->
        $scope.getStatus().then (result)->
          angular.forEach $scope.flows, (value, id)->
            if flow.id == parseInt(id, 10)
              $scope.flow.state = value
      , 500

    error = ->
      $state.go 'dashboard.flow.index'

    Flow.get_redactor {id: $stateParams.id}, success, error
    $scope.remove =->
      if confirm('точно удалить процесс?')
        remove()

  # get worker
  #------------------------------------------------------------------------------
  $scope.addNewWorker =->
    worker = new Worker({
      name: ''
      time: '* * * * * *'
      status: 'enabled'
      flow:
        id: parseInt($stateParams.id, 10)
      device_action:
        id: null
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

  # select2
  # ------------------
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

  # buttons remove|submit
  #------------------------------------------------------------------------------
  remove =->
    success =->
      $state.go("dashboard.flow.index")
    error =(result)->
      Message result.data.status, result.data.message
    $scope.flow.$delete success, error

  $scope.submit =->
    success =(data)->
      Notify 'success', 'Схема успешно сохранена', 3

    error =(result)->
      Message result.data.status, result.data.message

    scheme = $scope.callback.save()
    $scope.flow.objects = scheme.objects || []
    $scope.flow.connectors = scheme.connectors || []
    Flow.update_redactor {id: $stateParams.id}, $scope.flow, success, error

  # init
  #------------------------------------------------------------------------------
  getWorkflow()
  getFlow()

  vm
]