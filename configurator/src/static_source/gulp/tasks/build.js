var gulp = require('gulp'),
    runSequence = require('run-sequence');

gulp.task('default', function(cb) {
    runSequence(
        'build_lib_js',
        'build_coffee_js',
        ['build_haml'],
        ['build_templates'],
        'build_lib_css',
        'build_less',
        'webserver',
        'watch'
    );
});

gulp.task('pack', function(cb) {
    runSequence(
        'build_lib_js',
        'build_coffee_js',
        ['build_haml'],
        ['build_templates'],
        'build_lib_css',
        'build_less'
    );
});
