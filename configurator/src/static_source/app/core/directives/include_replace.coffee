#
# <div ng-include src="template" include-replace></div>
#

angular
.module('appDirectives')
.directive 'includeReplace', ["$log", ($log)->
  restrict: 'A'
  require: 'ngInclude'
  link: ($scope, element, attrs) ->
    element.replaceWith(element.children())
]
