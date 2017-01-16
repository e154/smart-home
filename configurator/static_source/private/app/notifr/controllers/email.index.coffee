angular
.module('appControllers')
.controller 'emailIndexCtrl', ['$scope', 'EmailTemplate','$state'
($scope, EmailTemplate, $state) ->

  vm = this

  tableCallback = {}
  vm.options =
    perPage: 20
    resource: EmailTemplate
    columns: [
      {
        name: 'notifr.system_name'
        field: 'name'
        width: '310px'
        template: '<strong>{{item[field]}}</strong>'
      }
      {
        name: 'notifr.description'
        field: 'description'
      }
      {
        name: 'notifr.created_at'
        field: 'created_at'
        width: '160px'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }
      {
        name: 'notifr.updated_at'
        field: 'updated_at'
        width: '160px'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }
    ]
    menu: null
    callback: tableCallback
    rowsClickCallback: ($event, item)->
      $event.preventDefault()
      $state.go('dashboard.notifr.template', {name: item.name})
      false
  vm
]
