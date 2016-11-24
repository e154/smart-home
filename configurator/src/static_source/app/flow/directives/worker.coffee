angular
.module('appDirectives')
.directive 'worker', ["$log", ($log)->
  restrict: 'A'
  scope:
    flow: '=flow'
  templateUrl: '/flow/templates/worker.html'
  link: ($scope, element, attrs) ->

]
