angular
.module('appControllers')
.controller 'roleIndexCtrl', ['$scope', 'Role', '$state', '$filter', 'Notify', 'Message'
  ($scope, Role, $state, $filter, Notify, Message) ->

    tableCallback = {}
    $scope.options =
      perPage: 100
      resource: Role
      columns: [
        {
          name: 'role.name'
          field: 'name'
          template: '<a ui-sref="dashboard.role.show({name: item[field]})"><strong>{{item[field]}}</strong></a>'
        }
        {
          name: 'role.parent'
          field: 'parent'
          template: '<span>{{item[field]["name"]}}</span>'
        }
        {
          name: 'role.description'
          field: 'description'
        }
        {
          name: 'role.created_at'
          field: 'created_at'
          template: '<span>{{item[field] | readableDateTime}}</span>'
          width: '120px'
        }
        {
          name: 'role.update_at'
          field: 'update_at'
          template: '<span>{{item[field] | readableDateTime}}</span>'
          width: '120px'
        }
      ]
      callback: tableCallback
      onLoad: (result)->
      rows: (item)->

]