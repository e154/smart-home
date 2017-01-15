#
# Created by delta54 on 12.03.15.
#
#
#
# <div class="table" table="ctrl.options"></div>
#
#  tableCallback = {}
#  tableCallback.update()
#  vm.options =
#    onLoad: (result)->
#    perPage: 20
#    resource: LayersliderResource
#    columns: [
#      {
#        name: 'id'
#        field: 'id'
#        width: '50px'
#      }
#      {
#        name: 'изображение'
#        field: 'data'
#        sort: 'disable'
#        template: "<div class='ls-preview-img' ng-init='image = angular.fromJson(item).layers[0].properties.background'>
#                                <a ng-href='#/layerslider/{{item.id}}'>
#                                    <img style='width: 100%' ng-src='#{vm.image_path}/{{item.id}}/{{column.getImage(item)}}' alt='' title='small'/>
#                                </a>
#                            </div>"
#        getImage: (slider)->
#          data = angular.fromJson(slider.data)
#          data.layers[0].properties.background
#
#        width: '250px'
#      }
#      {
#        name: 'name'
#        field: 'name'
#        width: '48%'
#      }
#      {
#        name: 'status'
#        field: 'status'
#        width: '100px'
#      }
#      {
#        name: 'create time'
#        field: 'create_time'
#        template: '<span>{{item[field] | date:"H:mm dd.MM.yyyy"}}</span>'
#      }
#      {
#        name: 'update time'
#        field: 'update_time'
#        template: '<span>{{item[field] | date:"H:mm dd.MM.yyyy"}}</span>'
#      }
#    ]
#    menu:
#      column: 2
#      buttons: [
#        {
#          name: 'edit'
#          clickCallback: ($event, item)->
#            $event.preventDefault()
#            $location.path "/layerslider/#{item.id}"
#            false
#        }
#        {
#          name: 'to trash'
#          showIf: (item)->
#            item.status != 'trash'
#
#          clickCallback: ($event, item)->
#            $event.preventDefault()
#            moveToTrash(item)
#            false
#        }
#        {
#          name: 'untrash'
#          showIf: (item)->
#            item.status == 'trash'
#
#          clickCallback: ($event, item)->
#            $event.preventDefault()
#            moveToTrash(item)
#            false
#        }
#        {
#          name: 'remove'
#          showIf: (item)->
#            item.status == 'trash'
#
#          clickCallback: ($event, item)->
#            $event.preventDefault()
#            remove(item)
#            false
#        }
#        {
#          name: 'copy'
#          clickCallback: ($event, item)->
#            $event.preventDefault()
#            copy(item)
#            false
#        }
#        {
#          name: 'export'
#          template: '<a target="_self" href="/api/admin/layerslider/{{item.id}}/export" download="">export</span>'
#        }
#      ]
#    rows: (item)->
#      color = ''
#      if item.status == 0
#        color = 'bg-success'
#      color
#    callback: tableCallback

#
#
#all:
#  method: 'POST'
#  responseType: 'json'
#  transformResponse: (data) ->
#    meta: data.meta
#    items: data.posts
#
#

'use strict'
angular
.module('appDirectives')
.directive 'table', ['$log','Message'
  ($log, Message)->
    restrict: 'A'
    template: '
<div class="clearfix" ng-if="pagination.objects_count >= perPage">
    <div class="pull-left">
        <pagination boundary-links="true" total-items="pagination.objects_count" items-per-page="perPage"
max-size="maxSize" ng-model="currentPage" class="pagination-sm" previous-text="&lsaquo;" next-text="&rsaquo;"
first-text="&laquo;" last-text="&raquo;"></pagination>
    </div>

    <div class="pull-right">
        <select class="form-control" name="amount" ng-model="perPage" ng-options="item for item in itemsPerPage"
ng-change="tableUpdate()"></select>
    </div>
</div>

<div ng-if="!items.length || !table.columns.length">
  Данные отсутствуют
</div>

<div ng-if="items.length && table.columns.length" class="form-group clearfix">{{"total items" | translate}} ({{pagination.objects_count}})</div>

<table ng-if="items.length && table.columns.length" class="table list-table with-menu">

<!--head-->
<thead>
<tr class="sortable">
    <th ng-repeat="column in table.columns" width="{{column.width}}">
        <a href="" ng-click="sortBy(column)" ng-if="column.sort != \'disable\'"
ng-class="{ \'sorted desc\': order == \'desc\' && sortby.indexOf(column.field) != -1, \'sorted asc\': order == \'asc\' && sortby.indexOf(column.field) != -1 }">
            <span>{{column.name | translate}}</span>
            <span class="sorting-indicator"></span>
        </a>
        <div ng-if="column.sort == \'disable\'"><span>{{column.name | translate}}</span></span></div>
    </th>
</tr>
</thead>
<!--/head-->

<!--body-->
<tbody>
<tr ng-repeat="item in items" ng-class="[table.rows(item)]" ng-style="{\'height\': table.menu ? \'57px\' : \'inherit\'}" ng-click="table.rowsClickCallback($event, item)">
  <td ng-repeat="column in table.columns" width="{{column.width}}">
    <a ng-if="column.clickCallback" href="" ng-click="column.clickCallback($event, item)">
      <span ng-if="!column.template">{{item[column.field || column.name]}}</span>
      <item-cell ng-if="column.template" item="item" field="column.field || column.name" column="column"></item-cell>
    </a>

    <span ng-if="!column.template && !column.clickCallback">{{item[column.field || column.name]}}</span>
    <item-cell ng-if="column.template && !column.clickCallback" item="item" field="column.field || column.name" column="column"></item-cell>

    <div class="row-actions" ng-if="$index == table.menu.column && table.menu.buttons">
      <span ng-repeat="button in table.menu.buttons" ng-if="showButton(button, item)">
        <span ng-if="!$first">|</span>
        <a ng-if="!button.template" href="" ng-click="button.clickCallback($event, item)">{{button.name}}</a>
        <button-cell ng-if="button.template" item="item" button="button" field="column.field || column.name" column="column"></button-cell>
      </span>
    </div>

  </td>

</tr>
</tbody>
<!--/body-->

</table>
'
    scope:
      table: "="

    link: ($scope, $element, attrs) ->

      pageRecalc =(perPage)->
        $scope.perPage = perPage
        $scope.itemsPerPage = [
          1 * $scope.perPage
          2 * $scope.perPage
          3 * $scope.perPage
          5 * $scope.perPage
          6 * $scope.perPage
        ]

      $scope.pagination =
        limit: 0
        objects_count: 0
        offset: 0
      pageRecalc(10)
      $scope.maxSize = 4
      $scope.currentPage = 1
      $scope.items = []
      $scope.query = {}
      $scope.sortby = ['created_at']
      $scope.order = ['desc']

      $scope.$watch 'table', (table)->
        return if !table
        if table.perPage?
          pageRecalc(table.perPage)

      # callback
      $scope.table.callback.query = (query)->
        $scope.query = query

      $scope.table.callback.update = ()->
        getItems()

      $scope.tableUpdate = getItems =->

        if !$scope.table.resource
          return

        success =(result)->
          $scope.items = angular.copy(result.items)
          $scope.pagination = angular.copy(result.meta)
          try
            $scope.table?.onLoad(result)
          catch

        error =(result)->
          Message result.data.status, result.data.message

        request =
          query: $scope.query
          sortby: $scope.sortby
          order: $scope.order
          limit: $scope.perPage
          offset: ($scope.currentPage - 1) * $scope.perPage

        if $scope.table.query
          request.query = $scope.table.query

        $scope.table.resource.all request, success, error

      $scope.sortBy =(c)->

        if $scope.sortby.indexOf(c.field) == -1
          $scope.sortby = []
          $scope.order = []
          $scope.sortby.push c.field
          $scope.order.push 'desc'
        else
          i = $scope.sortby.indexOf(c.field)
          if $scope.order[i] == 'asc'
            $scope.order[i] = 'desc'
          else
            $scope.order[i] = 'asc'

        getItems()

      $scope.showButton =(button, item)->
        if typeof button.showIf == 'undefined'
          return true
        else
          return button.showIf(item)

      $scope.$watch 'table + currentPage', ->
        getItems()

]

angular
.module('app')
  .directive 'itemCell', ['$log','$compile'
    ($log, $compile)->
      restrict: 'EA'
      scope:
        item: "=item"
        column: "=column"
        field: "=field"
      link: ($scope, $element, attrs) ->
        if $scope.column?.template
          template = $compile($scope.column.template)($scope)
          $element.append(template)

  ]

  .directive 'buttonCell', ['$log','$compile'
    ($log, $compile)->
      restrict: 'EA'
      scope:
        item: "=item"
        button: "=button"
        column: "=column"
        field: "=field"
      link: ($scope, $element, attrs) ->
        if $scope.button?.template
          console.log $scope.button.template
          $element.html($compile($scope.button.template)($scope))

  ]