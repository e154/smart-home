###
# Created by delta54 on 11.11.14.
###

angular
.module('appServices')
.factory 'storage', ['$log', ($log) ->
  class storage
    prefix: ""
    isOk: null

    constructor: (pref)->
      if pref
        @prefix = pref

      @isOk = @supports_html5_storage()

      this

    setKey: (name)->
      @prefix + '_' + name

    setItem: (key, value)->
      localStorage.setItem(@setKey(key), value.toString())

    setObject: (key, value)->
      @setItem(key, angular.toJson(value))

    getItem: (key)->
      localStorage.getItem(@setKey(key))

    getInt: (key)->
      if !@isOk
        return 0

      parseInt(@getItem(key),10)

    getFloat: (key)->
      if !@isOk
        return 1.0

      parseFloat(@getItem(key))

    getBool: (key)->
      if !@isOk
        return false

      return @getItem(key) == 'true'

    getObject: (key)->
      if !@isOk
        return {}

      return angular.fromJson(@getItem(key))

    supports_html5_storage: ()->
      return Storage?

    removeItem: (key)->
      localStorage.removeItem(@setKey(key));

  storage
]