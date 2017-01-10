angular
.module('appControllers')
.controller 'userIndexCtrl', ['$scope', 'User', '$state', '$filter'
($scope, User, $state, $filter) ->

  tableCallback = {}
  $scope.options =
    perPage: 100
    resource: User
    columns: [
      {
        name: 'user.nickname'
        field: 'nickname'
        template: '<strong>{{item[field]}}</strong>'
      }
      {
        name: 'user.status'
        field: 'status'
        width: '70px'
      }
      {
        name: 'user.role'
        field: 'role'
        template: '<span>{{item[field]["name"]}}</span>'
        width: '100px'
      }
      {
        name: 'user.first_name'
        field: 'first_name'
        width: '120px'
      }
      {
        name: 'user.email'
        field: 'email'
        width: '120px'
      }
      {
        name: 'user.created_at'
        field: 'created_at'
        template: '<span>{{item[field] | readableDateTime}}</span>'
        width: '120px'
      }
      {
        name: 'user.update_at'
        field: 'update_at'
        template: '<span>{{item[field] | readableDateTime}}</span>'
        width: '120px'
      }
    ]
    menu:
      column: 0
      buttons: [
        {
          name: $filter('translate')('user.menu.show')
          clickCallback: ($event, user)->
            $event.preventDefault()
            $state.go('dashboard.user.show', {id: user.id})
            false
        }
        {
          name: $filter('translate')('user.menu.edit')
          clickCallback: ($event, user)->
            $event.preventDefault()
            $state.go('dashboard.user.edit', {id: user.id})
            false
        }
        {
          name: $filter('translate')('user.menu.block')
          showIf: (user)->
            user.status != 'blocked'

          clickCallback: ($event, user)->
            $event.preventDefault()
            changeStatus(user, 'blocked')
            false
        }
        {
          name: 'unblock'
          showIf: (user)->
            user.status == 'blocked'

          clickCallback: ($event, user)->
            $event.preventDefault()
            changeStatus(user, 'active')
            false
        }
        {
          name: 'remove'
          showIf: (user)->
            user.status == 'blocked'

          clickCallback: ($event, user)->
            $event.preventDefault()
            deleteUser(user.id)
            false
        }
      ]
    callback: tableCallback
    onLoad: (result)->
    rows: (item)->

  deleteUser =(user)->
    console.log 'delete user', user

  changeStatus =(user)->
    console.log 'change status', user
]