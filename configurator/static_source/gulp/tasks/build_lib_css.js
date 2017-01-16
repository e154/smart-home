var gulp = require('gulp'),
    conf = require('../config').build_lib_css,
    mainBowerFiles = require('main-bower-files'),
    concat = require('gulp-concat'),
    filter = require('gulp-filter'),
    csso = require('gulp-csso'),
    inject = require("gulp-inject");

gulp.task('build_lib_css', function() {
    var cssFilter = filter('**/*.css');  //отбираем только css файлы
    return gulp.src(mainBowerFiles({
        includeDev: false,
        paths: conf.paths
    }))

        .pipe(cssFilter)
        //.pipe(csso()) // минимизируем css
        .pipe(concat(conf.filename))
        .pipe(gulp.dest(conf.dest)); // записываем css
});
