angular
.module('appControllers')
.controller 'mapIndexCtrl', ['$scope', 'MapResource', '$state', '$timeout', '$httpParamSerializer'
($scope, MapResource, $state, $timeout, $httpParamSerializer) ->

  tableCallback = {}
  $scope.options =
    perPage: 100
    resource: MapResource
    columns: [
      {
        name: '#'
        field: 'id'
        width: '50px'
      }
      {
        name: 'map.name'
        field: 'name'
        clickCallback: ($event, item)->
          $event.preventDefault()
          $state.go('dashboard.map.show', {id: item.id})
          false
      }
      {
        name: 'map.created_at'
        field: 'created_at'
        width: '140px'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }
      {
        name: 'map.update_at'
        field: 'update_at'
        width: '140px'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }


    ]
    menu: null
    callback: tableCallback
    onLoad: (result)->
    rows: (item)->

]