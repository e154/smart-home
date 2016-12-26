#
# example:
# #main-menu.light{'perfect-scrollbar'=>true}
# #main-menu.light{'perfect-scrollbar'=>'options'}
#

angular
.module('appServices')
.directive 'perfectScrollbar', ['$log', ($log) ->
  restrict: 'A'
  scope:
    perfectScrollbar: "="
  link: ($scope, $element, $attrs) ->

    options = {}

    if $scope.perfectScrollbar && typeof $scope.perfectScrollbar == 'Object'
      $.extend(options, $scope.perfectScrollbar)

    $element.perfectScrollbar(options)
]
