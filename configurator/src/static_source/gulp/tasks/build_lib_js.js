var gulp = require('gulp'),
    conf = require('../config').build_lib_js,
    mainBowerFiles = require('main-bower-files'),
    concat = require('gulp-concat'),
    filter = require('gulp-filter'),
    uglify = require('gulp-uglify'),
    inject = require("gulp-inject");

gulp.task('build_lib_js', function() {
    var jsFilter = filter('**/*.js');       //отбираем только  javascript файлы
    return gulp.src(mainBowerFiles({
        includeDev: false,
        paths: conf.paths
    }))

        .pipe(jsFilter)
        .pipe(concat(conf.filename))
        //.pipe(uglify())
        .pipe(gulp.dest(conf.dest));
});
