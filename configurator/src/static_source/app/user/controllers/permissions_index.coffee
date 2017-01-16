angular
.module('appControllers')
.controller 'permissionsIndexCtrl', ['$scope', '$stateParams', 'Role', 'Auth', 'Notify', 'Message', '$q'
($scope, $stateParams, Role, Auth, Notify, Message, $q) ->

  $scope.roles = new Role {
    limit:100
    offset:0
    order:'desc'
    query:{}
    sortby:'created_at'
  }
  $scope.roles.$all()
  $scope.access_list = new Auth {}
  $scope.getAccessList =->
    success =()->
    error =(result)->
      Message result.data.status, result.data.message
    $scope.access_list.$show success, error
  $scope.getAccessList()
]