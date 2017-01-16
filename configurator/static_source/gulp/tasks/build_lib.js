var gulp = require('gulp'),
    conf_css = require('../config').build_lib_css,
    conf_js = require('../config').build_lib_js,
    mainBowerFiles = require('main-bower-files'),
    concat = require('gulp-concat'),
    filter = require('gulp-filter'),
    csso = require('gulp-csso'),
    inject = require("gulp-inject"),
    uglify = require('gulp-uglify');

for (var key in conf_js) {

    (function (_task, _paths, _filename, _dest) {

        var task = 'bower:js:' + _task,
            paths = _paths,
            filename = _filename,
            dest = _dest;

        gulp.task(task, function() {
            var jsFilter = filter('**/*.js');       //отбираем только  javascript файлы
            return gulp.src(mainBowerFiles({
                includeDev: false,
                paths: paths
            }))

                .pipe(jsFilter)
                .pipe(concat(filename))
                .pipe(gulp.dest(dest));
        });

    })(key, conf_js[key].paths, conf_js[key].filename, conf_js[key].dest)
}

for (var key in conf_css) {

    (function (_task, _paths, _filename, _dest) {

        var task = 'bower:css:' + _task,
            paths = _paths,
            filename = _filename,
            dest = _dest;

        gulp.task(task, function () {
            var cssFilter = filter('**/*.css');  //отбираем только css файлы
            return gulp.src(mainBowerFiles({
                includeDev: false,
                paths: paths
            }))

                .pipe(cssFilter)
                //.pipe(csso()) // минимизируем css
                .pipe(concat(filename))
                .pipe(gulp.dest(dest)); // записываем css
        });

    })(key, conf_css[key].paths, conf_css[key].filename, conf_css[key].dest)
}

gulp.task('bower:js', ['bower:js:public', 'bower:js:private']);
gulp.task('bower:css', ['bower:css:public', 'bower:css:private']);
