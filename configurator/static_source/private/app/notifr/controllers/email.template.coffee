angular
.module('appControllers')
.controller 'emailTemplateCtrl', ['$scope','Message', 'EmailTemplate','$http', 'EmailItem', '$stateParams', '$state', 'Notify'
($scope, Message, EmailTemplate, $http, EmailItem, $stateParams, $state, Notify) ->

  vm = this

  vm.template =
    content: ''
    description: ''
    name: ''
  vm.template_name = $stateParams.name
  vm.items = []
  $scope.selectedItems = []
  vm.selectSize = 0
  vm.isNew = true
  vm.content =
    fields: []
    items: []
    title: ''

  getFrame = ->
    $('#preview').contents().find('html')

  vm.getTemplate = ->
    success =(template)->
      vm.isNew = false
      vm.template = template
      vm.content = angular.fromJson(template.content)
      $scope.selectedItems = vm.content.items
      # preview in iframe
      vm.preview()

    error =(response)->
      Message response.data.status, response.data.message

    EmailTemplate.show {name: vm.template_name}, success, error

  vm.getItems = ->
    success =(data)->
      vm.items = data.items
      vm.selectSize = vm.items.length

    error =(response)->
      Message response.data.status, response.data.message

    EmailItem.items {}, success, error

  vm.preview = ->
    if typeof vm.content == 'undefined'
      return
    vm.content.items = $scope.selectedItems
    $http.post(window.app_settings.server_url + '/api/v1/email/preview', vm.content)

    .success (data) ->
      frame = getFrame().html(data)
      frame.find('.edit_inline').editable()

    .error (data) ->
      Message data.status, data.message

  $scope.$watch 'selectedItems', () ->
    vm.preview()
    return

  serialize = ->
    array = []
    getFrame().find('.edit_inline').each (index) ->
      array.push
        name: $(this).attr('data-id')
        value: $(this).html()
      return
    vm.content.fields = array
    vm.content.items = $scope.selectedItems
    vm.template.content = angular.toJson(vm.content)
    return

  vm.save = ->
    serialize()

    success =(template)->
      Notify 'success', 'Шаблон успешно создан', 3
      vm.isNew = false
      vm.template_name = vm.template.name
      $state.go 'dashboard.notifr.template', {name: template.name}

    error =(response)->
      Message response.data.status, response.data.message

    EmailTemplate.create {}, vm.template, success, error

  vm.update = ->
    serialize()

    success =->
      Notify 'success', 'Шаблон успешно обновлен', 3
      vm.getTemplate()

    error =(response)->
      Message response.data.status, response.data.message

    EmailTemplate.update {name: vm.template.name}, vm.template, success, error

  vm.delete = ->
    success =->
      vm.isNew = true
      vm.template =
        content: ''
        description: ''
        name: ''
      vm.content =
        fields: []
        items: []
        title: ''
      vm.template_name = ''
      $scope.selectedItems = []
      getFrame().html ''
      $state.go 'dashboard.notifr.index', {}

    error =(response)->
      Message response.data.status, response.data.message

    if confirm('точно удалить шаблон?')
      EmailTemplate.delete {name: vm.template.name}, success, error


  #init
  vm.getItems()
  if vm.template_name != undefined
    vm.getTemplate()

  vm
]
