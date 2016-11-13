var gulp = require('gulp'),
    conf = require('../config').build_coffee_js,
    concat = require('gulp-concat'),
    ngconcat = require('gulp-ngconcat'),
    inject = require('gulp-inject'),
    coffee = require('gulp-coffee'),
    uglify = require('gulp-uglify'),
    ngClassify = require('gulp-ng-classify'),
    replace = require('gulp-replace');

var date = new Date();

gulp.task('build_coffee_js', function(done) {
    return gulp.src(conf.source)
        .pipe(coffee({bare: true})
            .on('error', done))     // Compile coffeescript
        //.pipe(ngconcat(conf.filename))
        .pipe(concat(conf.filename))
        //.pipe(uglify())
        //.pipe(ngClassify())
        .pipe(replace('__CURRENT_TIME__', date))
        .pipe(gulp.dest(conf.dest));

});
