#
# example:
# %file-manager(ng-model="model" multiple)
# %div(file-manager="true" ng-model="model")
#

angular
.module('appServices')
.directive 'fileManager', ['FileManager', '$compile', (FileManager, $compile) ->
  restrict: 'EA'
  transclude: true
  replace: true
  scope:
    fileManager: "="
    ngModel: "="
  template: '<div ng-click="onHandleClick($event)" ng-transclude></div>'
  link: ($scope, $element, $attrs) ->

    $scope.onHandleClick =->
      FileManager.show().then (images)=>
        if $attrs.multiple
          $scope.ngModel = images
        else
          $scope.ngModel = images[0]

]
