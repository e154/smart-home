var gulp = require('gulp'),
    connect = require('gulp-connect'),
    config = require('../config').webserver;

gulp.task('webserver', function() {
    connect.server({
        root: config.root,
        livereload: config.livereload,
        port: config.port,
        fallback: config.fallback
    });

    console.log('Server listening on http://localhost:'+config.port);
});

gulp.task('webserver_reload', function () {
    gulp.src(config.watch)
        .pipe(connect.reload());
});