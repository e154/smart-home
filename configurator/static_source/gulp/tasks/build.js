var gulp = require('gulp'),
    runSequence = require('run-sequence');

gulp.task('default', function(cb) {
    runSequence(
        'bower:css',
        'bower:js',
        'coffee',
        ['haml'],
        ['template'],
        'less',
        'copy',
        'redactor_theme_files', 'redactor_theme_less',
        'watch'
    );
});

gulp.task('pack', function(cb) {
    runSequence(
        'bower:css',
        'bower:js',
        'coffee',
        ['haml'],
        ['template'],
        'less',
        ['redactor_theme_files', 'redactor_theme_less'],
        'copy'
    );
});
