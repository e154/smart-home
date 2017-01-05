angular
.module('angular-map')
.directive 'mapViewerDevice', ['$timeout', 'Notify', 'Stream'
($timeout, Notify, Stream) ->
  restrict: 'A'
  replace: true
  templateUrl: '/map-viewer/templates/map_viewer_device.html'
  scope:
    element: '=mapViewerDevice'
  link: ($scope, $element, $attrs) ->

    # vars
    # --------------------
#    console.log 'element',$scope.element
    timer = null
    $scope.show_menu = false
    $scope.element.prototype.current_state = null
    st = $scope.element.graph_settings
    $element.css
      width: st.width || 'auto'
      height: st.height || 'auto'

    # stream
    # --------------------
    setState =(_state)->
      return if !_state || !_state?.id
      for map_element_state in $scope.element.prototype.states
        if map_element_state.device_state.id == _state.id
          $scope.element.prototype.current_state = map_element_state
          break

    $scope.$on 'broadcast_device_state', (e, data)->
      return if !data || !data?.state
      if $scope.element.prototype.device.id == data.id
        setState data.state

    $scope.doAction =(action, $event)->
      Stream.sendRequest("do.action", {action_id: action.device_action.id, device_id: $scope.element.prototype.device.id}).then (result)->
        if !result.error
          Notify 'success', "Команда выполнена успешно", 3
        else
          Notify 'error', "Результат выполнения команды:\n\r #{result.error}", 3

    # menu
    # --------------------
    $scope.mouseLive =->
      return if !$scope.show_menu
      timer = $timeout ()->
        $scope.show_menu = false
        timer = null
      , 2000

    $scope.mouseOver =->
      return if timer == null
      $timeout.cancel(timer)

    $scope.showMenu =->
      $scope.show_menu = !$scope.show_menu
      timer = null if !$scope.show_menu

    # etc
    # --------------------

    return
]