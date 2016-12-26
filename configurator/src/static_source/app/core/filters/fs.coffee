#
# <img ng-src="{{'5f9ec9826b2e3932f71966e74bdfcb76.jpg' | fs}}"/>
# <img ng-src="{{'5f9ec9826b2e3932f71966e74bdfcb76.jpg' | fs: '/attach/file_storage'}}"/>
#

angular
.module('appFilters')
.filter 'fs', ['$log', ($log) ->
  (name, dir) ->
    if !dir
      dir = "/attach/file_storage"
    for i in [0..2]
      dir += "/" + name.substring(i*3, (i+1)*3)
    return dir + "/" + name
]