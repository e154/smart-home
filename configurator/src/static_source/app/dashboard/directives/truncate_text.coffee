###*
# Created by delta54 on 09.11.14.
###

# usage:
#<p truncate-text="100">some text</p>
angular
.module('appDirectives')
.directive 'truncateText', ->
  restrict: 'A'
  link: ($scope, element, attrs) ->
# http://stackoverflow.com/questions/17138868/how-to-trigger-a-directive-when-updating-a-model-in-angularjs
    $scope.$watch 'someValue', (value) ->
      len = parseInt(attrs.truncateText)
      text = element.text().trim()
      trunc = undefined
      if text.length > len
        trunc = text.substring(0, len) + ' ...'
      element.text trunc
      return
    return
