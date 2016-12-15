###*
# Created by delta54 on 18.11.14.
###

angular.module('ui.bootstrap.pagination', [
  'template/pagination/pager.html'
  'template/pagination/pagination.html'
]).controller('PaginationController', [
  '$scope'
  '$attrs'
  '$parse'
  ($scope, $attrs, $parse) ->
    self = this
    ngModelCtrl = $setViewValue: angular.noop
    setNumPages = if $attrs.numPages then $parse($attrs.numPages).assign else angular.noop

    @init = (ngModelCtrl_, config) ->
      ngModelCtrl = ngModelCtrl_
      @config = config

      ngModelCtrl.$render = ->
        self.render()
        return

      if $attrs.itemsPerPage
        $scope.$parent.$watch $parse($attrs.itemsPerPage), (value) ->
          self.itemsPerPage = parseInt(value, 10)
          $scope.totalPages = self.calculateTotalPages()
          return
      else
        @itemsPerPage = config.itemsPerPage
      return

    @calculateTotalPages = ->
      totalPages = if @itemsPerPage < 1 then 1 else Math.ceil($scope.totalItems / @itemsPerPage)
      Math.max totalPages or 0, 1

    @render = ->
      $scope.page = parseInt(ngModelCtrl.$viewValue, 10) or 1
      return

    $scope.selectPage = (page) ->
      if $scope.page != page and page > 0 and page <= $scope.totalPages
        ngModelCtrl.$setViewValue page
        ngModelCtrl.$render()
      return

    $scope.getText = (key) ->
      $scope[key + 'Text'] or self.config[key + 'Text']

    $scope.noPrevious = ->
      $scope.page == 1

    $scope.noNext = ->
      $scope.page == $scope.totalPages

    $scope.$watch 'totalItems', ->
      $scope.totalPages = self.calculateTotalPages()
      return
    $scope.$watch 'totalPages', (value) ->
      setNumPages $scope.$parent, value
      # Readonly variable
      if $scope.page > value
        $scope.selectPage value
      else
        ngModelCtrl.$render()
      return
    return
]).constant('paginationConfig',
  itemsPerPage: 10
  boundaryLinks: false
  directionLinks: true
  firstText: 'First'
  previousText: 'Previous'
  nextText: 'Next'
  lastText: 'Last'
  rotate: true).directive('pagination', [
    '$parse'
    'paginationConfig'
    ($parse, paginationConfig) ->
      {
      restrict: 'EA'
      scope:
        totalItems: '='
        firstText: '@'
        previousText: '@'
        nextText: '@'
        lastText: '@'
      require: [
        'pagination'
        '?ngModel'
      ]
      controller: 'PaginationController'
      templateUrl: 'template/pagination/pagination.html'
      replace: true
      link: (scope, element, attrs, ctrls) ->
        paginationCtrl = ctrls[0]
        ngModelCtrl = ctrls[1]
        # Create page object used in template

        makePage = (number, text, isActive) ->
          {
          number: number
          text: text
          active: isActive
          }

        getPages = (currentPage, totalPages) ->
          pages = []
          # Default page limits
          startPage = 1
          endPage = totalPages
          isMaxSized = angular.isDefined(maxSize) and maxSize < totalPages
          # recompute if maxSize
          if isMaxSized
            if rotate
# Current page is displayed in the middle of the visible ones
              startPage = Math.max(currentPage - Math.floor(maxSize / 2), 1)
              endPage = startPage + maxSize - 1
              # Adjust if limit is exceeded
              if endPage > totalPages
                endPage = totalPages
                startPage = endPage - maxSize + 1
            else
# Visible pages are paginated with maxSize
              startPage = (Math.ceil(currentPage / maxSize) - 1) * maxSize + 1
              # Adjust last page if limit is exceeded
              endPage = Math.min(startPage + maxSize - 1, totalPages)
          # Add page number links
          number = startPage
          while number <= endPage
            page = makePage(number, number, number == currentPage)
            pages.push page
            number++
          # Add links to move between page sets
          if isMaxSized and !rotate
            if startPage > 1
              previousPageSet = makePage(startPage - 1, '...', false)
              pages.unshift previousPageSet
            if endPage < totalPages
              nextPageSet = makePage(endPage + 1, '...', false)
              pages.push nextPageSet
          pages

        if !ngModelCtrl
          return
        # do nothing if no ng-model
        # Setup configuration parameters
        maxSize = if angular.isDefined(attrs.maxSize) then scope.$parent.$eval(attrs.maxSize) else paginationConfig.maxSize
        rotate = if angular.isDefined(attrs.rotate) then scope.$parent.$eval(attrs.rotate) else paginationConfig.rotate
        scope.boundaryLinks = if angular.isDefined(attrs.boundaryLinks) then scope.$parent.$eval(attrs.boundaryLinks) else paginationConfig.boundaryLinks
        scope.directionLinks = if angular.isDefined(attrs.directionLinks) then scope.$parent.$eval(attrs.directionLinks) else paginationConfig.directionLinks
        paginationCtrl.init ngModelCtrl, paginationConfig
        if attrs.maxSize
          scope.$parent.$watch $parse(attrs.maxSize), (value) ->
            maxSize = parseInt(value, 10)
            paginationCtrl.render()
            return
        originalRender = paginationCtrl.render

        paginationCtrl.render = ->
          originalRender()
          if scope.page > 0 and scope.page <= scope.totalPages
            scope.pages = getPages(scope.page, scope.totalPages)
          return

        return

      }
  ]).constant('pagerConfig',
  itemsPerPage: 10
  previousText: '« Previous'
  nextText: 'Next »'
  align: true).directive 'pager', [
  'pagerConfig'
  (pagerConfig) ->
    {
    restrict: 'EA'
    scope:
      totalItems: '='
      previousText: '@'
      nextText: '@'
    require: [
      'pager'
      '?ngModel'
    ]
    controller: 'PaginationController'
    templateUrl: 'template/pagination/pager.html'
    replace: true
    link: (scope, element, attrs, ctrls) ->
      paginationCtrl = ctrls[0]
      ngModelCtrl = ctrls[1]
      if !ngModelCtrl
        return
      # do nothing if no ng-model
      scope.align = if angular.isDefined(attrs.align) then scope.$parent.$eval(attrs.align) else pagerConfig.align
      paginationCtrl.init ngModelCtrl, pagerConfig
      return

    }
]
angular.module('template/pagination/pager.html', []).run [
  '$templateCache'
  ($templateCache) ->
    $templateCache.put 'template/pagination/pager.html', '<ul class="pager">\n' + '  <li ng-class="{disabled: noPrevious(), previous: align}"><a href ng-click="selectPage(page - 1)">{{getText(\'previous\')}}</a></li>\n' + '  <li ng-class="{disabled: noNext(), next: align}"><a href ng-click="selectPage(page + 1)">{{getText(\'next\')}}</a></li>\n' + '</ul>'
    return
]
angular.module('template/pagination/pagination.html', []).run [
  '$templateCache'
  ($templateCache) ->
    $templateCache.put 'template/pagination/pagination.html', '<ul class="pagination">\n' + '  <li ng-if="boundaryLinks" ng-class="{disabled: noPrevious()}"><a href ng-click="selectPage(1)">{{getText(\'first\')}}</a></li>\n' + '  <li ng-if="directionLinks" ng-class="{disabled: noPrevious()}"><a href ng-click="selectPage(page - 1)">{{getText(\'previous\')}}</a></li>\n' + '  <li ng-repeat="page in pages track by $index" ng-class="{active: page.active}"><a href ng-click="selectPage(page.number)">{{page.text}}</a></li>\n' + '  <li ng-if="directionLinks" ng-class="{disabled: noNext()}"><a href ng-click="selectPage(page + 1)">{{getText(\'next\')}}</a></li>\n' + '  <li ng-if="boundaryLinks" ng-class="{disabled: noNext()}"><a href ng-click="selectPage(totalPages)">{{getText(\'last\')}}</a></li>\n' + '</ul>'
    return
]
