angular
.module('appControllers')
.controller 'userIndexCtrl', ['$scope', 'User', '$state', '$filter', 'Notify', 'Message'
($scope, User, $state, $filter, Notify, Message) ->

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
            updateStatus(user, 'blocked')
            false
        }
        {
          name: $filter('translate')('user.menu.unblock')
          showIf: (user)->
            user.status == 'blocked'

          clickCallback: ($event, user)->
            $event.preventDefault()
            updateStatus(user, 'active')
            false
        }
        {
          name: $filter('translate')('user.menu.remove')
          showIf: (user)->
            user.status == 'blocked'

          clickCallback: ($event, user)->
            $event.preventDefault()
            remove(user.id)
            false
        }
      ]
    callback: tableCallback
    onLoad: (result)->
    rows: (item)->

  remove =(id)->
    return if !confirm('точно удалить пользователя?')
    success =->
      tableCallback.update()
      return
    error =(result)->
      Message result.data.status, result.data.message
    User.delete {id: id}, success, error
    return

  updateStatus =(user, status)->
    return if !status || status == ''
    user.status = status
    success =->
      Notify 'success', 'Пользователь успешно обновлён', 3
    error =(result)->
      Message result.data.status, result.data.message
    User.update_status user, success, error
]