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

    options = $.extend(true, defaultOptions, $scope.options)

    $scope.onHandleClick =->
      FileManager.show(options).then (images)=>
        if options.multiple
          $scope.ngModel = images
        else
          $scope.ngModel = images[0]
]
