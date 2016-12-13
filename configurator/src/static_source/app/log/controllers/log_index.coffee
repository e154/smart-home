angular
.module('appControllers')
.controller 'logIndexCtrl', ['$scope', 'Log', '$state', '$timeout'
($scope, Log, $state, $timeout) ->

  tableCallback = {}
  $scope.options =
    perPage: 50
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

]