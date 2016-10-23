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
    gulp.watch(config.webserver.watch, ['webserver_reload']);
});
