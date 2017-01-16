var gulp = require('gulp'),
    conf = require('../config').build_haml,
    //connect = require('gulp-connect'),
    cache = require('gulp-cached'),
    haml = require('gulp-haml');

gulp.task('build_haml', function() {
    return gulp.src(conf.source, {read: true})
        .pipe(cache('linting'))
        .pipe(haml())
        .pipe(gulp.dest(conf.tmp));
        //.pipe(connect.reload());
});