(function($) {
    "use strict"; // Start of use strict

    //===============================================================
    // HIGHLIGHT
    //===============================================================
    $(document).ready(function() {
        $('pre code').each(function(i, block) {
            hljs.highlightBlock(block);
        });

        //fancybox
        $('.fancybox').fancybox();
    });

    //===============================================================
    // ETC
    //===============================================================

    // Disable empty links in docs examples
    $('[href="#"]').click(function (e) {
        e.preventDefault()
    });

    // jQuery for page scrolling feature - requires jQuery Easing plugin
    $('a.page-scroll').bind('click', function(event) {
        var $anchor = $(this);
        console.log($anchor.attr('href'));
        $('html, body').stop().animate({
            scrollTop: ($($anchor.attr('href')).offset().top - 0)
        }, 1250, 'easeInOutExpo');
        event.preventDefault();
    });

    $('body').scrollspy({
        target: '#docsNavbarContent',
        offset: 150
    });

    // Closes the Responsive Menu on Menu Item Click
    $('.navbar-collapse ul li a').click(function() {
        $('.navbar-toggle:visible').click();
    });

    // Offset for Main Navigation
    // $('#mainNav').affix({
    //     offset: {
    //         top: 100
    //     }
    // });

    $('#docsNavbarContent').affix({
        offset: {
            top: 200
        }
    });

})(jQuery); // End of use strict

(function () {
    'use strict';

    anchors.options.placement = 'left';
    anchors.add('.docs-section > h1, .docs-section > h2, .docs-section > h3, .docs-section > h4, .docs-section > h5')
}());
