var gulp = require('gulp'),
    config = require('../config');

gulp.task('watch', function() {

    //  ...
    //================//
    gulp.watch(config.build_coffee_js.watch, function() {
        gulp.run('build_coffee_js');
    });
    gulp.watch(config.build_less.watch, function() {
        gulp.run('build_less');
    });
    gulp.watch(config.build_haml.watch, function() {
        gulp.run('build_haml');
    });
    gulp.watch(config.build_templates.watch, function() {
        gulp.run('build_templates');
    });
    gulp.watch(config.redactor_theme_less.watch, function() {
        gulp.run('redactor_theme_less');
    });
    gulp.watch(config.redactor_theme_files.watch, function() {
        gulp.run('redactor_theme_files');
    });
});
