#
# example:
# {{file.size | size}}
#

angular
.module('appFilters')
.filter 'size', ['$log', ($log) ->
  (value, si)->
    if !value
      return ''

    thresh = if si then 1000 else 1024
    if Math.abs(bytes) < thresh
      return bytes + ' B'
    units = if si then [
      'kB','MB','GB','TB','PB','EB','ZB','YB'
    ] else [
      'KiB','MiB','GiB','TiB','PiB','EiB','ZiB','YiB',]
    u = -1
    loop
      bytes /= thresh
      ++u
      unless Math.abs(bytes) >= thresh and u < units.length - 1
        break
    bytes.toFixed(1) + ' ' + units[u]
]