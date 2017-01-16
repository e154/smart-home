angular
.module('appControllers')
.controller 'permissionsShowCtrl', ['$scope', '$stateParams', 'Role', 'Notify', 'Message', '$q'
($scope, $stateParams, Role, Notify, Message, $q) ->

  return if !$stateParams.name

  $scope.role = new Role {name: $stateParams.name}
  $scope.role_access_list = {}
  current_access_list = {}

  getRole =->
    success =()->
    error =(result)->
      Message result.data.status, result.data.message
    $scope.role.$show success, error

  getRoleAccessList =->
    success =(role_access_list)->
      $scope.role_access_list = role_access_list
    error =(result)->
      Message result.data.status, result.data.message
    Role.get_access_list {name: $stateParams.name}, success, error

  $scope.update =->
    success =->
      Notify 'success','Права доступа успешно обновлены.',3
      current_access_list = {}
    error =(result)->
      Message result.data.status, result.data.message
    Role.update_access_list {name: $stateParams.name}, current_access_list, success, error

  $scope.checked =(pack_name, level_name, level)->
    current_access_list[pack_name] = {} if !current_access_list[pack_name]
    current_access_list[pack_name][level_name] = level.checked || false

  $scope.cancel =->
    $scope.$parent.getAccessList()
    current_access_list = {}
    getRole()
    getRoleAccessList()

  # init
  getRole()
  getRoleAccessList()



]