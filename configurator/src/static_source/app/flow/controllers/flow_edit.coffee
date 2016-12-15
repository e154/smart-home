angular
.module('appControllers')
.controller 'flowEditCtrl', ['$scope', 'Message', '$stateParams', 'Flow', '$state', 'Workflow', '$timeout'
'log', 'Notify', 'Worker', '$rootScope'
($scope, Message, $stateParams, Flow, $state, Workflow, $timeout, log, Notify, Worker, $rootScope) ->
  vm = this

  # vars
  $scope.callback = {}
  $scope.workflows = []
  $scope.flow = {}
  $scope.elementScripts = {}
  $scope.elementFlows = {}

  # watcher
  #------------------------------------------------------------------------------
  instance = $rootScope.$on '$stateChangeStart', (event, toState, toParams, fromState, fromParams, options)->
    if !confirm('Вы точно хотите покинут редактирование процесса?')
      event.preventDefault()
    $scope.$on('$destroy', instance);

  #------------------------------------------------------------------------------
  # workflow list
  #------------------------------------------------------------------------------
  getWorkflow =->
    success = (result)->
      $scope.workflows = result.items
    error = (result)->
      Message result.data.status, result.data.message
    Workflow.all {}, success, error

  #------------------------------------------------------------------------------
  # flow
  #------------------------------------------------------------------------------
  getFlow =->
    success = (flow) ->
      $scope.flow = flow
      if !$scope.flow?.workers
        $scope.flow.workers = []

      # scripts
      angular.forEach $scope.flow.objects, (object)->
        $scope.elementScripts[object.id] = object.script if object.script?.id?
        $scope.elementFlows[object.id] = object.flow_link if object.flow_link?

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

  #------------------------------------------------------------------------------
  # remove
  #------------------------------------------------------------------------------
  remove =->
    $scope.$on('$destroy', instance);
    success =->
      $state.go("dashboard.flow.index")
    error =(result)->
      Message result.data.status, result.data.message
    $scope.flow.$delete success, error

  #------------------------------------------------------------------------------
  # save
  #------------------------------------------------------------------------------
  $scope.submit =->
    success =(data)->
      $scope.$on('$destroy', instance);
      Notify 'success', 'Схема успешно сохранена', 3

    error =(result)->
      Message result.data.status, result.data.message

    scheme = $scope.callback.save()
    $scope.flow.objects = scheme.objects || []

    # scripts & flows
    angular.forEach $scope.flow.objects, (object)->
      object.script = $scope.elementScripts[object.id] || null
      object.flow_link = $scope.elementFlows[object.id] || null

    $scope.flow.connectors = scheme.connectors || []
    Flow.update_redactor {id: $stateParams.id}, $scope.flow, success, error

  #------------------------------------------------------------------------------
  # init
  #------------------------------------------------------------------------------
  getWorkflow()
  getFlow()

  vm
]