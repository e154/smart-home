#
# example:
# %file-manager(ng-model="model" multiple)
# %div(file-manager="{}" ng-model="model")
#

angular
.module('appServices')
.directive 'fileManager', ['FileManager', '$compile', (FileManager, $compile) ->
  restrict: 'EA'
  transclude: true
  replace: true
  scope:
    options: "=fileManager"
    ngModel: "="
  template: '<div ng-click="onHandleClick($event)" ng-transclude></div>'
  link: ($scope, $element, $attrs) ->

    $scope.options = {} if !$scope.options

    defaultOptions =
      multiple: false

    $scope.onHandleClick =->
      FileManager.show($.extend(true, defaultOptions, options)).then (images)=>
        $scope.ngModel = images
]
