angular
.module('appServices')
.service 'authForm', [ '$rootScope', 'ngDialog'
($rootScope, ngDialog) ->
  class authForm

    ngDialogInstance: null
    scope: null

    constructor: ->
      @scope = $rootScope.$new()
      @scope.$on 'event:auth-loginRequired', (event, data) =>
        @show()
      return

    show: ->
      return if @ngDialogInstance
      @ngDialogInstance = ngDialog.open
        template: '/core/templates/loginform.html'
        controller: 'loginFormCtrl'
        className: 'ngdialog-theme-default dashboard-login-form'
        scope: @scope
        plain: false
        overlay: true
        showClose: false
        preCloseCallback: ()=>
          @ngDialogInstance = null
          return

      return

  new authForm()
]