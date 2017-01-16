var gulp = require('gulp'),
    conf = require('../config').build_less,
    less = require('gulp-less'),
    concat = require('gulp-concat'),
    gutil = require('gulp-util');

gulp.task('build_less', function(done) {
    return gulp.src(conf.source)
        .pipe(concat(conf.filename))
        .pipe(less())
        .on('error', function(err){
            gutil.log(err);
            this.emit('end');
        })
        .pipe(gulp.dest(conf.dest));
});