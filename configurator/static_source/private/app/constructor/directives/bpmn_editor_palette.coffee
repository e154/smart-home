angular
.module('angular-bpmn')
.directive 'bpmnEditorPalette', ['log', (log) ->
    restrict: 'A'
    scope:
      bpmnEditorPalette: '='
      settings: '='
    template: '<div class="group" ng-repeat="group in bpmnEditorPalette.groups" data-group="{{::group.name}}">
<div class="entry" ng-repeat="entry in group.items" ng-class="[entry.class]" bpmn-editor-palette-node="entry" entry-type="{{entry.type}}" settings="settings" data-help="{{entry.tooltip | translate}}"></div>'
    link: ($scope, element, attrs) ->
  ]
