angular
.module('appControllers')
.controller 'flowEditCtrl', ['$scope', 'Message', '$stateParams', 'Flow', '$state', 'Workflow', '$timeout', 'log'
($scope, Message, $stateParams, Flow, $state, Workflow, $timeout, log) ->
  vm = this

  # vars
  $scope.callback = {}
  $scope.workflows = []

  # workflow list
  #------------------------------------------------------------------------------
  success = (result)->
    console.log result
    $scope.workflows = result.items
  error = (result)->
    Message result.data.status, result.data.message
  Workflow.all {}, success, error

  # get flow
  #------------------------------------------------------------------------------
  $scope.flow = {}
  success = (flow) ->
    $scope.flow = flow
    $timeout ()->
      $scope.getStatus().then (result)->
        $scope.flows = result.flows

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
#      $state.go("dashboard.flow.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    scheme = $scope.callback.save()
    $scope.flow.objects = scheme.objects || []
    $scope.flow.connectors = scheme.connectors || []
    Flow.update_redactor {id: $stateParams.id}, $scope.flow, success, error

  vm
]