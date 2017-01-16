angular
.module('appControllers')
.controller 'logIndexCtrl', ['$scope', 'Log', '$state', '$timeout', '$httpParamSerializer'
($scope, Log, $state, $timeout, $httpParamSerializer) ->

  $scope.levels = [
    {
      name: 'Emergency'
      value: 'Emergency'
      selected: false
    }
    {
      name: 'Alert'
      value: 'Alert'
      selected: false
    }
    {
      name: 'Critical'
      value: 'Critical'
      selected: false
    }
    {
      name: 'Error'
      value: 'Error'
      selected: false
    }
    {
      name: 'Warning'
      value: 'Warning'
      selected: false
    }
    {
      name: 'Notice'
      value: 'Notice'
      selected: false
    }
    {
      name: 'Info'
      value: 'Info'
      selected: false
    }
    {
      name: 'Debug'
      value: 'Debug'
      selected: false
    }
  ]
  $scope.start_date
  $scope.end_date

  tableCallback = {}
  $scope.options =
    perPage: 100
    resource: Log
    columns: [
      {
        name: 'log.created_at'
        field: 'created_at'
        width: '140px'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }
      {
        name: 'log.level'
        field: 'level'
        width: '70px'
      }
      {
        name: 'log.body'
        field: 'body'
      }
    ]
    menu: null
    callback: tableCallback
    onLoad: (result)->
    rows: (item)->
      style
      switch item.level
        when 'Emergency'
          style = 'log-emergency'
        when 'Alert'
          style = 'log-alert'
        when 'Critical'
          style = 'log-critical'
        when 'Error'
          style = 'log-error'
        when 'Warning'
          style = 'log-warning'
        when 'Notice'
          style = 'log-notice'
        when 'Info'
          style = 'log-info'
        when 'Debug'
          style = 'log-debug'
      style

  $scope.$watch 'levels', (val)->
    return if !val
    $scope.update()

  , true

  $scope.setEndDay =(date)->
    $scope.end_date = date
    $scope.update()

  $scope.setStartDay =(date)->
    $scope.start_date = date
    $scope.update()

  $scope.update =->
    selected = []
    angular.forEach $scope.levels, (level)->
      if level.selected
        selected.push "'#{level.value}'"

    query = {}

    if $scope.start_date
      query.start_date = moment($scope.start_date).format("YYYY-MM-DD")

    if $scope.end_date
      query.end_date = moment($scope.end_date).format("YYYY-MM-DD")

    if selected.length
      query.levels = selected.join(',')

    tableCallback.query(query)
    tableCallback.update()

]