!(() ->
  !do ->
    if window.console and 'undefined' != typeof console.log
      url = 'https://github.com/e154/smart-home'
      i = "Modern smart systems for life - #{url}"
      console.log('%c Smart home %c Copyright © 2014-%s', 'font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;font-size:62px;color:#303E4D;-webkit-text-fill-color:#303E4D;-webkit-text-stroke: 1px #303E4D;', 'font-size:12px;color:#a9a9a9;', (new Date).getFullYear())
      console.log('%c ' + i, 'color:#333;')
)()