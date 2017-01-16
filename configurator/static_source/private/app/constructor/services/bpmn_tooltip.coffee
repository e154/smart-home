do ->
  'use strict'
  window.lsTooltip =
    timeout: 0
    init: ->
      $(document).on 'mouseover', '[data-help]', ->
        el = this
        lsTooltip.timeout = setTimeout((->
          lsTooltip.open el
          return
        ), 400)
        return
      $(document).on 'mouseout', '[data-help]', ->
        clearTimeout lsTooltip.timeout
        lsTooltip.close()
        return
      return
    destroy: ->
      $(document).off 'mouseover', '[data-help]'
      $(document).off 'mouseout', '[data-help]'
      return
    open: (el) ->
      # Create tooltip
      $('body').prepend $('<div>', 'class': 'ls-tooltip').append($('<div>', 'class': 'inner')).append($('<span>'))
      # Get tooltip
      tooltip = $('.ls-tooltip')
      # Set tooltip text
      tooltip.find('.inner').html $(el).data('help')
      # Get viewport dimensions
      v_w = $(window).width()
      # Get element dimensions
      e_w = $(el).width()
      # Get element position
      e_l = $(el).offset().left
      e_t = $(el).offset().top
      # Get toolip dimensions
      t_w = tooltip.outerWidth()
      t_h = tooltip.outerHeight()
      # Position tooltip
      tooltip.css
        top: e_t - t_h - 10
        left: e_l - ((t_w - e_w) / 2)
      # Fix right position
      if tooltip.offset().left + t_w > v_w
        tooltip.css
          'left': 'auto'
          'right': 10
        tooltip.find('span').css
          left: 'auto'
          right: v_w - ($(el).offset().left) - ($(el).outerWidth() / 2) - 17
          marginLeft: 'auto'
      return
    close: ->
      $('.ls-tooltip').remove()
      return

  $(document).ready ->
    window.lsTooltip.init()
    return
  return
