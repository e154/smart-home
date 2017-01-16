var gulp = require('gulp'),
    conf_files = require('../config').redactor_theme_files,
    conf_theme = require('../config').redactor_theme_less,
    less = require('gulp-less'),
    concat = require('gulp-concat'),
    gutil = require('gulp-util');

gulp.task('redactor_theme_files', function() {
    return gulp.src(conf_files.minimal.source)
        .pipe(gulp.dest(conf_files.minimal.dest));
});

gulp.task('redactor_theme_less', function() {
    return gulp.src(conf_theme.minimal.source)
        .pipe(concat(conf_theme.filename))
        .pipe(less())
        .on('error', function(err){
            gutil.log(err);
            this.emit('end');
        })
        .pipe(gulp.dest(conf_theme.minimal.dest));
});
