$(document).ready(function () {
    $.each($('.navbar-nav').find('li').find('a'), function() {
        if ($(this).attr('id') === window.location.pathname) {
            $(this).addClass('active')
        }
    });
});