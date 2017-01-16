angular
.module('appControllers')
.controller 'deviceStatusCtrl', ['$scope', 'Notify', 'DeviceState', 'Message', '$stateParams', '$timeout'
($scope, Notify, DeviceState, Message, $stateParams, $timeout) ->
  vm = this

  vm.statuses = []
  vm.getAll =->
    query =
      limit:100
      offset:0
      order:'desc'
      query:
        device_id: $stateParams.id
      sortby:'created_at'

    success =(result)->
      vm.statuses = angular.copy(result.items) || []

    error =(result)->
      Message result.data.status, result.data.message

    DeviceState.all query, success, error

  vm.addNew =->
    vm.statuses.push new DeviceState({
      system_name: "NAME"
      description: ""
      device:
        id: parseInt($stateParams.id, 10)
    })

  vm.remove =(_state)->
    return if !confirm('Вы точно хотите удалить это состояние?')

    if _state.id
      success =()->
        statuses = angular.copy(vm.statuses)
        angular.forEach statuses, (state, key)->
          if _state.system_name == state.system_name
            statuses.splice(key, 1)
        vm.statuses = angular.copy(statuses)

      error =(result)->
        Message result.data.status, result.data.message

      state = new DeviceState _state
      state.$delete success, error

    else
      statuses = angular.copy(vm.statuses)
      angular.forEach statuses, (state, key)->
        if _state.system_name == state.system_name
          statuses.splice(key, 1)
      vm.statuses = angular.copy(statuses)

  vm.update =(_state)->
    success =(result)->
      angular.forEach vm.statuses, (status, key)->
        if status.system_name == result.system_name
          vm.statuses[key] = angular.copy result

    error =(result)->
      Message result.data.status, result.data.message

    state = new DeviceState _state
    state.$update success, error

  vm.create =(_state)->
    success =(result)->
      angular.forEach vm.statuses, (status, key)->
        if status.system_name == result.system_name
          vm.statuses[key] = angular.copy result

    error =(result)->
      Message result.data.status, result.data.message

    state = new DeviceState _state
    state.$create success, error

  vm.getAll()

  vm
]