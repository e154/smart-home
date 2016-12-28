angular
.module('appDirectives')
.directive 'fancybox', ["log", (log)->
  restrict: 'A'
  scope:
    options: '=fancybox'
  link: ($scope, element, attrs) ->
    if !$scope.fancybox
      options = {}

    defaultOptions =
      helpers:
        title :
          type : ''
        overlay :
          speedOut : 0
          locked: false

    element.fancybox($.extend true, defaultOptions, options)

]
