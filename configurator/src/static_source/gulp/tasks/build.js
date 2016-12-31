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
        'ace_themes',
        'redactor_theme_files', 'redactor_theme_less',
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
        'build_less',
        ['redactor_theme_files', 'redactor_theme_less'],
        'ace_themes'
    );
});
