var gulp = require('gulp'),
    conf = require('../config').redactor_theme_files,
    less = require('gulp-less'),
    concat = require('gulp-concat'),
    gutil = require('gulp-util');

gulp.task('redactor_theme_files', function() {
    return gulp.src(conf.minimal.source)
        .pipe(gulp.dest(conf.minimal.dest));
});
